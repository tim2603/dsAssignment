// distribute tasks, heartbeat + masterelection, hearbeats from workers -> deactivation,
// does master give out tasks or do workers ask for tasks

package main

import (
	"context"
	"ds/grpc/ds"
	"ds/grpc/general"
	logging "ds/grpc/logger"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logger = logging.Logger{}

type Server struct {
	ds.UnsafeCommunicationWithMasterServiceServer
	master *Master
}

type State int64

const (
	Alive State = iota
	Dead
)

type Master struct {
	mapTaskStatus                map[string]int
	reduceTaskStatus             map[int]int
	finish                       bool
	inputFiles                   []string
	nReduce                      int
	taskPackageAmount            int
	taskPackageRest              int
	intermediateFiles            []string
	mapFinished                  bool
	reduceFinished               bool
	MutexLock                    sync.Mutex
	activeWorkers                []WorkerStatus
	interval                     int
	task                         *general.MapReduceTask
	currentTournamentTreeRecords map[string]RecordEntry
	sortedTournamentTreeRecords  []RecordEntry
	threshholdForRecordsToWrite  int
}

type WorkerStatus struct {
	timestamp    timestamppb.Timestamp
	workerID     string // workerID = IP address => Is it okay?
	workerState  State
	taskState    general.TaskState
	workerClient ds.CommunicationWithWorkerServiceClient
	taskSlice    []string
}

type RecordEntry struct {
	record string
	value  float32
}

// SendCurrentTournamentTreeValue implements ds.CommunicationWithMasterServiceServer.
func (server *Server) SendCurrentTournamentTreeValue(context context.Context, value *ds.TournamentTreeValue) (*ds.EmptyMessage, error) {
	server.master.onReceivedCurrentTournamentTreeValue(value.GetValue(), value.GetWorkerId().WorkerId)
	return &ds.EmptyMessage{}, nil
}

func (m *Master) getValueFromRecordEntry(entry string) float32 {
	value, error := strconv.ParseFloat(strings.Split(entry, " ")[3], 32)
	if error != nil {
		logger.Debug(error.Error())
	}
	return float32(value)
}

func (m *Master) onReceivedCurrentTournamentTreeValue(entry string, workerID string) {
	logger.Debug("Received entry " + entry + " from worker " + workerID)
	m.MutexLock.Lock()
	m.currentTournamentTreeRecords[workerID] = RecordEntry{record: entry, value: m.getValueFromRecordEntry(entry)}
	if len(m.currentTournamentTreeRecords) == m.task.N_reducers {
		// TODO: Das muss eigentlich die Anzahl der intermediate files sein

		logger.Debug("Beginning tournament tree iteration")
		// find smallest value
		var smallestValue RecordEntry = RecordEntry{record: "", value: math.MaxFloat32}
		var workerWithSmallestValue string = ""
		for workerID_i, entry_i := range m.currentTournamentTreeRecords {
			value_i := entry_i.value
			if value_i < smallestValue.value {
				smallestValue = entry_i
				workerWithSmallestValue = workerID_i
			}
		}

		// var smallestValue float32 = math.MaxFloat32
		// var workerWithSmallestValue string = ""
		// for workerID_i, value_i := range m.currentTournamentTreeRecords {
		// 	if value_i < smallestValue {
		// 		smallestValue = value_i
		// 		workerWithSmallestValue = workerID_i
		// 	}
		// }

		delete(m.currentTournamentTreeRecords, workerWithSmallestValue)

		if len(m.sortedTournamentTreeRecords) == m.threshholdForRecordsToWrite {
			m.writeSortedRecordsToFile(m.sortedTournamentTreeRecords)
			m.sortedTournamentTreeRecords = make([]RecordEntry, 0, m.threshholdForRecordsToWrite)
		}
		m.sortedTournamentTreeRecords = append(m.sortedTournamentTreeRecords, smallestValue)

		logger.Debug("Selecting value " + fmt.Sprintf("%f", smallestValue.value) + " from worker " + workerWithSmallestValue)

		for i := range m.activeWorkers {
			if m.activeWorkers[i].workerID == workerWithSmallestValue {
				go m.activeWorkers[i].workerClient.SendSelectedTournamentTreeSelection(context.Background(), &ds.TournamentTreeSelection{Selected: true})
				// 	go func(i int) {
				// 		m.activeWorkers[i].workerClient.SendSelectedTournamentTreeSelection(context.Background(), &ds.TournamentTreeSelection{Selected: true})
				// 	}(i)
				// } else {
				// 	go func(i int) {
				// 		m.activeWorkers[i].workerClient.SendSelectedTournamentTreeSelection(context.Background(), &ds.TournamentTreeSelection{Selected: false})
				// 	}(i)
			}
		}
	}
	m.MutexLock.Unlock()

}

func (m *Master) writeSortedRecordsToFile(records []RecordEntry) {
	newFile, err := os.OpenFile("FinalFile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	}
	defer newFile.Close()

	for i := range records {
		newFile.WriteString(records[i].record + "\n")
	}

	newFile.Sync()
}

// if workers are not responding, set the status and reshedule the task
func (m *Master) CheckWorkerStatus() {
	m.MutexLock.Lock()
	currentTime := timestamppb.Now().AsTime()
	for index, s := range m.activeWorkers {
		if currentTime.Sub(m.activeWorkers[index].timestamp.AsTime()) > time.Duration(m.interval)*time.Second {
			if m.activeWorkers[index].workerState == Alive {
				m.activeWorkers[index].workerState = Dead
				logger.Debug("Worker got killed: " + s.workerID)
				//TODO: stop that task and give to another (only kinda for map phase)

				if m.activeWorkers[index].taskState == general.Mapping {
					m.redistributeMapTask(m.activeWorkers[index].taskSlice)
				}
				if m.activeWorkers[index].taskState == general.Reducing {
					m.redistributeReduceTask(m.activeWorkers[index].taskSlice)
					m.task.N_reducers = m.task.N_reducers - 1
				}
			}
		} //else {
		/*if m.activeWorkers[index].workerState == Dead {
			if m.task.State == general.Mapping {
				m.activeWorkers[index].taskState = general.AfterMapping
			}
			if m.activeWorkers[index].taskState == general.Reducing {
				m.activeWorkers[index].taskState = general.AfterReducing
			}
		}*/
		//s.workerState = Alive

		//}
	}
	m.MutexLock.Unlock()

	//TODO: maybe delete
	if m.task.State == general.BeforeMapping {
		if m.getCountOfActiveAndFullyConnectedWorkers() == m.task.N_mappers {
			m.startMapPhase()
		}
	} else if m.task.State == general.BeforeReducing {
		if m.getCountOfActiveAndFullyConnectedWorkers() == m.task.N_reducers {
			m.startReducePhase()
		}
	}
}

func (m *Master) redistributeMapTask(fileSlice []string) {
	/*for i := range m.activeWorkers {
		if (m.activeWorkers[i].workerState == Alive) && (m.activeWorkers[i].taskState == general.AfterMapping) {
			m.activeWorkers[i].taskState = general.Mapping
			m.activeWorkers[i].taskSlice = append(m.activeWorkers[i].taskSlice, fileSlice...)
			m.activeWorkers[i].workerClient.AssignMapTask(context.Background(), &ds.MapTask{Filenames: fileSlice})
		}
	}*/

	for i := range m.activeWorkers {
		if m.activeWorkers[i].workerState == Alive {
			m.activeWorkers[i].taskSlice = append(m.activeWorkers[i].taskSlice, fileSlice...)
			m.activeWorkers[i].taskState = general.Mapping
			m.activeWorkers[i].workerClient.AssignMapTask(context.Background(), &ds.MapTask{Filenames: m.activeWorkers[i].taskSlice})
			return
		}
	}
}

func (m *Master) redistributeReduceTask(fileSlice []string) {
	logger.Debug("Redistributing Reduce task")
	/*for i := range m.activeWorkers {
		if (m.activeWorkers[i].workerState == Alive) && (m.activeWorkers[i].taskState == general.AfterMapping) {
			m.activeWorkers[i].taskState = general.Mapping
			m.activeWorkers[i].taskSlice = append(m.activeWorkers[i].taskSlice, fileSlice...)
			m.activeWorkers[i].workerClient.AssignMapTask(context.Background(), &ds.MapTask{Filenames: fileSlice})
		}
	}*/
	for i := range m.activeWorkers {
		if m.activeWorkers[i].workerState == Alive {
			go m.activeWorkers[i].workerClient.RestartTournamentTree(context.Background(), &ds.EmptyMessage{})
		}
	}
	for i := range m.activeWorkers {
		if m.activeWorkers[i].workerState == Alive {
			m.activeWorkers[i].taskSlice = append(m.activeWorkers[i].taskSlice, fileSlice...)
			m.activeWorkers[i].taskState = general.Mapping
			go m.activeWorkers[i].workerClient.AssignReduceTask(context.Background(), &ds.ReduceTask{IntermediateFile: m.activeWorkers[i].taskSlice})

			m.currentTournamentTreeRecords = make(map[string]RecordEntry, m.getCountOfActiveAndFullyConnectedWorkers())
			m.sortedTournamentTreeRecords = make([]RecordEntry, 0, 10000)
			return
		}
	}
}

func (m *Master) getCountOfActiveAndFullyConnectedWorkers() int {
	count := 0
	for i := range m.activeWorkers {
		if m.activeWorkers[i].workerState == Alive && m.activeWorkers[i].workerClient != nil {
			count++
		}
	}
	return count
}

func (m *Master) startMapPhase() {
	m.task.State = general.Mapping
	logger.Debug("Starting Map phase with " + strconv.Itoa(len(m.activeWorkers)) + " workers")
	m.distributeCalculation(len(m.inputFiles), general.BeforeMapping)
	a := 0
	b := 0
	counter := m.taskPackageRest
	length := len(m.inputFiles)

	// test := m.inputFiles

	for i := range m.activeWorkers {
		if (m.activeWorkers[i].taskState == general.BeforeMapping) && (m.activeWorkers[i].workerState == Alive) {

			a = i * m.taskPackageAmount
			b = (i + 1) * m.taskPackageAmount
			fileSlice := make([]string, b-a)
			copy(fileSlice, m.inputFiles[a:b])
			// fileSlice := m.inputFiles[a:b]
			fmt.Println(counter)
			if counter != 0 {
				fileSlice = append(fileSlice, m.inputFiles[(length)-counter:(length)-(counter-1)]...)
				counter--
			}
			m.activeWorkers[i].taskSlice = fileSlice
			m.activeWorkers[i].taskState = general.Mapping
			go func(i int) {
				m.activeWorkers[i].workerClient.AssignMapTask(context.Background(), &ds.MapTask{Filenames: fileSlice})
			}(i)
		}
	}
}

// func (m *Master) distributeCalculation() {
// 	readyWorkers := 0
// 	for i := range m.activeWorkers {
// 		if (m.activeWorkers[i].workerState == Alive) && (m.activeWorkers[i].taskState == general.BeforeMapping) {
// 			readyWorkers = readyWorkers + 1
// 		}
// 	}
// 	m.taskPackageRest = len(m.inputFiles) % int(readyWorkers)
// 	m.taskPackageAmount = len(m.inputFiles) / int(readyWorkers)
// }

func (m *Master) distributeCalculation(packages int, taskState general.TaskState) {
	readyWorkers := 0
	for i := range m.activeWorkers {
		if (m.activeWorkers[i].workerState == Alive) && (m.activeWorkers[i].taskState == taskState) {
			readyWorkers = readyWorkers + 1
		}
	}
	fmt.Println(readyWorkers)
	m.taskPackageRest = packages % int(readyWorkers)
	m.taskPackageAmount = packages / int(readyWorkers)
}

func (m *Master) startReducePhase() {
	// TODO: Intermediate File discovery

	// TODO: copy map phase
	// intermediate files auch auf google
	// wait for every answer --> const communication via tournament tree method
	m.task.State = general.Reducing
	m.intermediateFiles = listFilesInDir("../worker/intermediate-files/")
	m.distributeCalculation(len(m.intermediateFiles), general.AfterMapping)
	logger.Debug("Starting Reduce phase with " + strconv.Itoa(len(m.activeWorkers)) + " workers")
	a := 0
	b := 0
	for i := range m.activeWorkers {
		if (m.activeWorkers[i].taskState == general.AfterMapping) && (m.activeWorkers[i].workerState == Alive) {
			m.activeWorkers[i].taskState = general.Reducing
			a = i * m.taskPackageAmount
			b = (i + 1) * m.taskPackageAmount
			fmt.Println(a)
			fmt.Println(b)
			fmt.Println(m.taskPackageAmount)
			logger.Debug("Intermediate files: " + strings.Join(m.intermediateFiles, ", "))
			fileSlice := m.intermediateFiles[a:b]

			m.activeWorkers[i].taskSlice = fileSlice
			go func(i int) {
				m.activeWorkers[i].workerClient.AssignReduceTask(context.Background(), &ds.ReduceTask{IntermediateFile: fileSlice})
			}(i)
		}
	}

}

// Called when receiving a heartbeat from Worker
func (s *Server) SendHeartBeatToMaster(context context.Context, heartbeat *ds.Heartbeat) (*ds.EmptyMessage, error) {
	s.master.OnReceivedHeartbeatFromWorker(heartbeat.GetWorkerId().GetWorkerId())
	return &ds.EmptyMessage{}, nil
}

// Called when notified about a finished map task from Worker
func (s *Server) NotifyAboutFinishedMapTask(context context.Context, workerID *ds.WorkerID) (*ds.EmptyMessage, error) {
	s.master.OnNotificationAboutFinishedMapTask(workerID.GetWorkerId())
	return &ds.EmptyMessage{}, nil
}

// Called when notified about a finished reduce task from Worker
func (s *Server) NotifyAboutFinishedReduceTask(context context.Context, workerID *ds.WorkerID) (*ds.EmptyMessage, error) {
	s.master.OnNotificationAboutFinishedReduceTask(workerID.GetWorkerId())
	return &ds.EmptyMessage{}, nil
}

// Called when receiving an address from a worker
func (s *Server) SendAddressToMaster(context context.Context, address *ds.Address) (*ds.EmptyMessage, error) {
	s.master.onReceivedAddressFromWorker(address.GetAddress(), address.GetWorkerId().GetWorkerId())
	return &ds.EmptyMessage{}, nil
}

func (m *Master) onReceivedAddressFromWorker(address string, workerId string) {
	logger.Debug("Received address " + address + " from worker with ID " + workerId)
	m.OnReceivedHeartbeatFromWorker(workerId) // Hacky
	m.MutexLock.Lock()
	defer m.MutexLock.Unlock()
	for i := range m.activeWorkers {
		if m.activeWorkers[i].workerID == workerId {
			m.activeWorkers[i].workerClient = getWorkerClient(address, workerId)
			logger.Debug("Assigned address " + address + " to worker with ID " + workerId)
			break
		}
	}
	// If enough workers are connected, start map phase
	if m.task.State == general.BeforeMapping {
		if m.getCountOfActiveAndFullyConnectedWorkers() == m.task.N_mappers {
			m.startMapPhase()
		}
	}
	// else if m.task.State == general.BeforeReducing {
	// 	if m.getCountOfActiveAndFullyConnectedWorkers() == m.task.N_reducers {
	// 		m.startReducePhase()
	// 	}
	// }
}

func (m *Master) OnReceivedHeartbeatFromWorker(workerID string) {
	m.MutexLock.Lock()
	logger.Debug("Received heartbeat from worker with ID " + workerID)
	currentTime := timestamppb.Now()

	isInList := false
	for index, s := range m.activeWorkers {
		if s.workerID == workerID {
			// s.timestamp = *currentTime
			m.activeWorkers[index].timestamp = *currentTime
			// fmt.Println(currentTime.AsTime().String() + "In der Schleife")
			isInList = true
		}
	}
	if !isInList {
		m.activeWorkers = append(m.activeWorkers, WorkerStatus{timestamp: *currentTime, workerID: workerID, workerState: Alive, taskState: general.BeforeMapping})
	}
	m.MutexLock.Unlock()
}

func (m *Master) OnNotificationAboutFinishedMapTask(workerID string) {
	logger.Debug("Map task for file " + workerID + " finished")
	for i, worker := range m.activeWorkers {
		if worker.workerID == workerID {
			m.activeWorkers[i].taskState = general.AfterMapping
		}

	}
	if m.areAllWorkerDoneWithMapping() {
		m.task.State = general.BeforeReducing
		// m.intermediateFiles = listFilesInDir("../worker/intermediate_files/")
		m.startReducePhase()
	}
}

func (m *Master) areAllWorkerDoneWithMapping() bool {
	for _, worker := range m.activeWorkers {
		if worker.workerState == Alive {
			// print(worker.workerID + " ")
			// println(worker.workerState)
			if worker.taskState != general.AfterMapping {
				return false
			}
		}
	}
	return true
}

func (m *Master) OnNotificationAboutFinishedReduceTask(workerID string) {
	m.task.N_reducers--
	logger.Debug("Reduce task for file " + workerID + " finished")
	for i, worker := range m.activeWorkers {
		if worker.workerID == workerID {
			m.activeWorkers[i].taskState = general.AfterReducing
		}

	}
	if m.areAllWorkerDoneWithReducing() {
		m.task.State = general.AfterReducing
		m.writeSortedRecordsToFile(m.sortedTournamentTreeRecords)
		logger.Debug("MapReduce-task finished!")
		fmt.Println("MapReduce-task finished!")
	}
}
func (m *Master) areAllWorkerDoneWithReducing() bool {
	for _, worker := range m.activeWorkers {
		if worker.workerState == Alive {
			if worker.taskState != general.AfterReducing {
				return false
			}
		}
	}
	return true
}

func getWorkerClient(address string, workerId string) ds.CommunicationWithWorkerServiceClient {
	serverAddr := flag.String(workerId, address, "The server address in the format of host:port")
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Debug(err.Error())
		return nil
	}
	client := ds.NewCommunicationWithWorkerServiceClient(conn)
	logger.Debug("Connected with client server at " + address)
	return client
}

func (master *Master) StartWorkerCheck() {
	ticker := time.NewTicker(20 * time.Second)
	time.Sleep(20 * time.Second)
	// Run the periodic function in a goroutine
	go func(master *Master) {
		for {
			select {
			case <-ticker.C:
				master.CheckWorkerStatus()
			}
		}
	}(master)

	// Google cloud Bucket / storage
	// What runs well on what conditions
	// Strategy at presentation, nothing beforehand
	// Different amount of mappers and reducers -> Intermediate files (hash for each word -> mod operation to assign it to a intermediate file)
	// maybe channels

	// Run the main program for a long, long, long seconds
	time.Sleep(time.Until(time.Date(2024, time.December, time.Now().Day(), 0, 0, 0, 0, time.Local)))

	// Stop the ticker when done
	ticker.Stop()
	logger.Debug("MasterCheck stopped. Program exiting.")
}

func listFilesInDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Error(err.Error())
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = dir + entry.Name()
	}
	logger.Debug("Files in directory: " + strings.Join(names, ", "))
	fmt.Println("Files in directory: " + strings.Join(names, ", "))
	return names
}

func main() {
	// file := initLogger()

	logger.Init("../../logs/master.log")
	logger.Debug("Starting")
	n_workers, error := strconv.Atoi(os.Args[1])
	if error != nil {
		logger.Error(error.Error())
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	filenames := listFilesInDir("../../../data/")

	master := &Master{interval: 15, task: &general.MapReduceTask{N_mappers: n_workers, N_reducers: n_workers}, inputFiles: filenames[:], currentTournamentTreeRecords: make(map[string]RecordEntry, 6), sortedTournamentTreeRecords: make([]RecordEntry, 0, 10000), threshholdForRecordsToWrite: 10000}
	ds.RegisterCommunicationWithMasterServiceServer(s, &Server{master: master})
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("Failed to serve: %v", err)
		}
	}()

	master.StartWorkerCheck()

	defer logger.CloseFile()
}
