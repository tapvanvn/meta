
file=cluster.jsonc
num_master_nodes=$(go run ..//main.go $file length master.nodes)
num_master_index=$(expr $num_master_nodes - 1)
echo "num master nodes :$num_master_nodes"
for i in $(seq 0 $num_master_index)
do 
    cmd="go run ../main.go  $file value master.nodes.$i"
    
    hostname=$($cmd);

    cmd2="go run ../main.go  $file value master.nodes.$hostname"

    mac=$($cmd2); 

    echo "#$i $hostname $mac"
    echo "\t$cmd\n\t$cmd2"
done

#go run main/main.go key "master.nodes.0"
#go run main/main.go value "master.nodes.master-1"