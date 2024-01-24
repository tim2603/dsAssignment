package main

import (
	"bufio"
	"context"
	"ds/grpc/ds"
	logging "ds/grpc/logger"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = logging.Logger{}

type Worker struct {
	masterID     string // Maybe rename it to workerID since it is used in the master to identify the worker
	masterClient ds.CommunicationWithMasterServiceClient

	smallest_value_pointer    int
	MutexLock                 sync.Mutex
	intermediate_file_content []string
	cloudStorageClient        *storage.Client
}

type Server struct {
	ds.UnsafeCommunicationWithWorkerServiceServer
	worker *Worker
}

func (worker *Worker) connectToCloudStorage() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("Connected to Google Cloud Storage")
	worker.cloudStorageClient = client
}

func (s *Server) RestartTournamentTree(context.Context, *ds.EmptyMessage) (*ds.EmptyMessage, error) {
	s.worker.MutexLock.Lock()
	s.worker.smallest_value_pointer = 0
	s.worker.MutexLock.Unlock()
	return &ds.EmptyMessage{}, nil
}

func (s *Server) AssignMapTask(context context.Context, mapTask *ds.MapTask) (*ds.EmptyMessage, error) {
	s.worker.onAssignedMapTask(mapTask)
	return &ds.EmptyMessage{}, nil
}

func (s *Server) AssignReduceTask(context context.Context, reduceTask *ds.ReduceTask) (*ds.EmptyMessage, error) {
	s.worker.onAssignedReduceTask(reduceTask)
	return &ds.EmptyMessage{}, nil
}

func (s *Server) SendSelectedTournamentTreeSelection(context context.Context, selection *ds.TournamentTreeSelection) (*ds.EmptyMessage, error) {
	s.worker.onReceivedTournamentTreeSelection(selection.GetSelected())
	return &ds.EmptyMessage{}, nil
}

func (worker *Worker) onAssignedMapTask(mapTask *ds.MapTask) {
	logger.Debug("Master assigned map-task " + strings.Join(mapTask.GetFilenames(), ", ") + " to me.")
	print("Master assigned map-task " + strings.Join(mapTask.GetFilenames(), ", ") + " to me.")
	worker.mapping(mapTask)
	worker.masterClient.NotifyAboutFinishedMapTask(context.Background(), &ds.WorkerID{WorkerId: worker.masterID})
}

func (worker *Worker) concatFiles(filenames []string) []string {
	var s []string
	s = make([]string, 0)
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)

		bkt := worker.cloudStorageClient.Bucket("distributed_systems2024")
		obj := bkt.Object(filename)
		reader, error := obj.NewReader(context.Background())
		if error != nil {
			logger.Error(error.Error())
		}
		scanner := bufio.NewScanner(reader)

		// scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			s = append(s, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
		}
	}
	return s
}

func (worker *Worker) mapping(mapTask *ds.MapTask) {
	// This is the place where the mapping happens
	logger.Debug("Mapping " + mapTask.String())
	//sort my stuff
	// The following code reads the data from the given file line by line
	// There has to be logic that takes the correct file

	concattedFiles := worker.concatFiles(mapTask.GetFilenames())
	worker.sortFiles(concattedFiles)

	//save in intermediate file
	workerId := worker.masterID

	worker.writeContentToFileInCloudStorage(concattedFiles, "intermediate-files/Intermediate-file-"+workerId+".txt")
	// newFile, err := os.Create("./intermediate-files/Intermediate-file-" + workerId)
	// if err != nil {
	// }
	// defer newFile.Close()

	// for _, line := range concattedFiles {
	// 	newFile.WriteString(line + "\n")
	// }

	// newFile.Sync()

	logger.Debug("Finished dummy map task")
}

func (worker *Worker) writeContentToFileInCloudStorage(content []string, filename string) {
	logger.Debug("Writing content to file " + filename + " in Cloud Storage")
	bkt := worker.cloudStorageClient.Bucket("distributed_systems2024")

	obj := bkt.Object(filename)
	newFile := obj.NewWriter(context.Background())
	// if err != nil {
	// }
	defer newFile.Close()

	for _, line := range content {
		newFile.Write([]byte(line + "\n"))
	}
}

func (worker *Worker) sortFiles(files []string) []string {
	sort.Slice(files, func(i, j int) bool {
		numberA, _ := strconv.ParseFloat(strings.Split(files[i], " ")[3], 32)
		numberB, _ := strconv.ParseFloat(strings.Split(files[j], " ")[3], 32)
		return numberA < numberB
	})
	return files
}

func (worker *Worker) reducing(reduceTask *ds.ReduceTask) {
	// This is the place where the reducing happens
	// TODO: Read intermediate file
	files := reduceTask.GetIntermediateFile()
	if len(files) > 1 {
		worker.intermediate_file_content = worker.sortFiles(worker.concatFiles(files))
	} else {
		worker.intermediate_file_content = worker.concatFiles(files)
	}

	logger.Debug("Reducing " + fmt.Sprintf("%d", len(worker.intermediate_file_content)) + " entries")
	for worker.smallest_value_pointer < len(worker.intermediate_file_content) {
		logger.Debug("Sending current tree entry of " + worker.intermediate_file_content[worker.smallest_value_pointer])
		worker.masterClient.SendCurrentTournamentTreeValue(context.Background(), &ds.TournamentTreeValue{Value: worker.intermediate_file_content[worker.smallest_value_pointer], WorkerId: &ds.WorkerID{WorkerId: worker.masterID}})
		time.Sleep(time.Second * 2)
	}
	logger.Debug("Finished reduce task")
	worker.masterClient.NotifyAboutFinishedReduceTask(context.Background(), &ds.WorkerID{WorkerId: worker.masterID})
}

func (worker *Worker) onReceivedTournamentTreeSelection(selected bool) {
	logger.Debug("Worker with ID " + worker.masterID + " selected!")
	worker.MutexLock.Lock()
	if selected {
		worker.smallest_value_pointer = worker.smallest_value_pointer + 1
	}
	worker.MutexLock.Unlock()
	if worker.smallest_value_pointer < len(worker.intermediate_file_content) {
		worker.masterClient.SendCurrentTournamentTreeValue(context.Background(), &ds.TournamentTreeValue{Value: worker.intermediate_file_content[worker.smallest_value_pointer], WorkerId: &ds.WorkerID{WorkerId: worker.masterID}})
	}

}

func (worker *Worker) onAssignedReduceTask(reduceTask *ds.ReduceTask) {
	logger.Debug("Master assigned reduce-task " + strings.Join(reduceTask.GetIntermediateFile(), ", ") + " to me.")
	worker.reducing(reduceTask)
}

func getMasterClient(address string) ds.CommunicationWithMasterServiceClient {
	serverAddr := flag.String("addr", address, "The server address in the format of host:port")
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	client := ds.NewCommunicationWithMasterServiceClient(conn)
	logger.Debug("Connected with master server at " + address)
	return client
}

func (worker *Worker) startHeartbeats() {
	ticker := time.NewTicker(5 * time.Second)

	// Run the periodic function in a goroutine
	go func() {

		for {
			select {
			case <-ticker.C:
				worker.masterClient.SendHeartBeatToMaster(context.Background(), &ds.Heartbeat{WorkerId: &ds.WorkerID{WorkerId: worker.masterID}})
				logger.Debug("Sending heartbeat")
			}
		}
	}()

	// Run the main program for a long, long, long seconds
	time.Sleep(time.Until(time.Date(2025, time.December, time.Now().Day(), 0, 0, 0, 0, time.Local)))

	// Stop the ticker when done
	ticker.Stop()
	logger.Debug("Ticker stopped. Program exiting.")
}

func main() {
	own_port := os.Args[3]
	master_port := os.Args[2]
	master_id := os.Args[1]
	master_address := "localhost:" + master_port
	own_address := "localhost:" + own_port
	logger.Init("../../logs/worker-" + master_id + ".log")
	logger.Debug("Starting")

	worker := &Worker{masterID: master_id, masterClient: getMasterClient(master_address)}
	worker.connectToCloudStorage()

	lis, err := net.Listen("tcp", ":"+own_port)
	if err != nil {
		logger.Error("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ds.RegisterCommunicationWithWorkerServiceServer(s, &Server{worker: worker})
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error("Failed to serve: %v", err)
		}
		defer s.Stop()
	}()
	worker.masterClient.SendAddressToMaster(context.Background(), &ds.Address{Address: own_address, WorkerId: &ds.WorkerID{WorkerId: worker.masterID}})
	worker.startHeartbeats()
	defer logger.CloseFile()
	// worker.startHeartbeats()
}
