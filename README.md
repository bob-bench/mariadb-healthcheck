# Health checking of MariaDB Galera on Azure LoadBalancer

Simple health checking for Mariadb Galera for the Microsoft Azure
LoadBalancer. Handle the HTTP queries of the Azure LB, check the
state of the cluster and provide feedback.

# TODO

[ ] Make options configurable
[ ] Make it reliable


# Current state
Note that development is currently halted as the MySQL driver fails basic
reliability testing. For the load balancer the following reliability tests
were considered:

1. Cluster not running (TCP SYN answered with RST)
2. Cluster IP ending in black hole (connection timeout)
3. Network trouble (no data, slow data)
4. Cluster very busy ("slowloris like apache issue")
4. System not a primary


The network trouble/cluster busy was simulated by running `nc -l 3306` and
besides the usage of a context.Context the backend would stall forever. As
there is no way to kill a go-routine it means that this software can't be
reliable yet.
