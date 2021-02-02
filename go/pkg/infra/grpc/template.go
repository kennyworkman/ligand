package grpc

const template = `# An unique identifier for the head node and workers of this cluster.
cluster_name: default

# The minimum number of workers nodes to launch in addition to the head
# node. This number should be >= 0.
min_workers: 0

# The maximum number of workers nodes to launch in addition to the head
# node. This takes precedence over min_workers.
max_workers: $ARG_MAX_WORKERS

# The autoscaler will scale up the cluster faster with higher upscaling speed.
# E.g., if the task requires adding more nodes then autoscaler will gradually
# scale up the cluster in chunks of upscaling_speed*currently_running_nodes.
# This number should be > 0.
upscaling_speed: 1.0

# This executes all commands on all nodes in the docker container,
# and opens all the necessary ports to support the Ray cluster.
# Empty string means disabled.
docker:
  image: "rayproject/ray-ml:latest-gpu" # You can change this to latest-cpu if you don't need GPU support and want a faster startup
  # image: rayproject/ray:latest-gpu   # use this one if you don't need ML dependencies, it's faster to pull
  container_name: "ray_container"
  # If true, pulls latest version of image. Otherwise, 'docker run' will only pull the image
  # if no cached version is present.
  pull_before_run: True
  run_options: []  # Extra options to pass into "docker run"

  # Example of running a GPU head with CPU workers
  # head_image: "rayproject/ray-ml:latest-gpu"
  # Allow Ray to automatically detect GPUs

  # worker_image: "rayproject/ray-ml:latest-cpu"
  # worker_run_options: []

# If a node is idle for this many minutes, it will be removed.
idle_timeout_minutes: 5

# Cloud-provider specific configuration.
provider:
  type: aws
  region: us-west-2
  # Availability zone(s), comma-separated, that nodes may be launched in.
  # Nodes are currently spread between zones by a round-robin approach,
  # however this implementation detail should not be relied upon.
  availability_zone: us-west-2a,us-west-2b
  # Whether to allow node reuse. If set to False, nodes will be terminated
  # instead of stopped.
  cache_stopped_nodes: False #https://github.com/ray-project/ray/issues/12104

# How Ray will authenticate with newly launched nodes.
auth:
  ssh_user: ubuntu
# By default Ray creates a new private keypair, but you can also use your own.
# If you do so, make sure to also set "KeyName" in the head and worker node
# configurations below.
#    ssh_private_key: /path/to/your/key.pem

# Provider-specific config for the head node, e.g. instance type. By default
# Ray will auto-configure unspecified fields such as SubnetId and KeyName.
# For more documentation on available fields, see:
# http://boto3.readthedocs.io/en/latest/reference/services/ec2.html#EC2.ServiceResource.create_instances
head_node:
  InstanceType: $ARG_INSTANCE
  ImageId: ami-0a2363a9cff180a64 # Deep Learning AMI (Ubuntu) Version 30

  # You can provision additional disk space with a conf as follows
  BlockDeviceMappings:
    - DeviceName: /dev/sda1
      Ebs:
        VolumeSize: 100

  # Additional options in the boto docs.

# Provider-specific config for worker nodes, e.g. instance type. By default
# Ray will auto-configure unspecified fields such as SubnetId and KeyName.
# For more documentation on available fields, see:
# http://boto3.readthedocs.io/en/latest/reference/services/ec2.html#EC2.ServiceResource.create_instances
worker_nodes:
  InstanceType: m5.large
  ImageId: ami-0a2363a9cff180a64 # Deep Learning AMI (Ubuntu) Version 30

  # Run workers on spot by default. Comment this out to use on-demand.
  InstanceMarketOptions:
    MarketType: spot
    # Additional options can be found in the boto docs, e.g.
    #   SpotOptions:
    #       MaxPrice: MAX_HOURLY_PRICE

  # Additional options in the boto docs.

# Files or directories to copy to the head and worker nodes. The format is a
# dictionary from REMOTE_PATH: LOCAL_PATH, e.g.
file_mounts: {
#    "/path1/on/remote/machine": "/path1/on/local/machine",
#    "/path2/on/remote/machine": "/path2/on/local/machine",
}

# Files or directories to copy from the head node to the worker nodes. The format is a
# list of paths. The same path on the head node will be copied to the worker node.
# This behavior is a subset of the file_mounts behavior. In the vast majority of cases
# you should just use file_mounts. Only use this if you know what you're doing!
cluster_synced_files: []

# Whether changes to directories in file_mounts or cluster_synced_files in the head node
# should sync to the worker node continuously
file_mounts_sync_continuously: False

# Patterns for files to exclude when running rsync up or rsync down
rsync_exclude:
  - "**/.git"
  - "**/.git/**"

# Pattern files to use for filtering out files when running rsync up or rsync down. The file is searched for
# in the source directory and recursively through all subdirectories. For example, if .gitignore is provided
# as a value, the behavior will match git's behavior for finding and using .gitignore files.
rsync_filter:
  - ".gitignore"

# List of commands that will be run before 'setup_commands'. If docker is
# enabled, these commands will run outside the container and before docker
# is setup.
initialization_commands: []

# List of shell commands to run to set up nodes.
setup_commands: 
$ARG_DEP_LIST
  # Note: if you're developing Ray, you probably want to create a Docker image that
  # has your Ray repo pre-cloned. Then, you can replace the pip installs
  # below with a git checkout <your_sha> (and possibly a recompile).
  # Uncomment the following line if you want to run the nightly version of ray (as opposed to the latest)
  # - pip install -U https://s3-us-west-2.amazonaws.com/ray-wheels/latest/ray-2.0.0.dev0-cp37-cp37m-manylinux2014_x86_64.whl

# Custom commands that will be run on the head node after common setup.
head_setup_commands: []

# Custom commands that will be run on worker nodes after common setup.
worker_setup_commands: []

# Command to start ray on the head node. You don't need to change this.
head_start_ray_commands:
  - ray stop
  - ulimit -n 65536; ray start --head --port=6379 --object-manager-port=8076 --autoscaling-config=~/ray_bootstrap_config.yaml

# Command to start ray on worker nodes. You don't need to change this.
worker_start_ray_commands:
  - ray sto
  - ulimit -n 65536; ray start --address=$RAY_HEAD_IP:6379 --object-manager-port=8076 `
