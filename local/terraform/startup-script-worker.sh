mkdir test123
sudo apt-get install git -y
git clone https://github.com/LUH-VSS/project-ds-PatrickMi
cd project-ds-PatrickMi/real_project/grpc/worker
./worker $(hostname) 50051 50061