# Health checking of MariaDB Galera on Azure LoadBalancer

Simple health checking for Mariadb Galera for the Microsoft Azure
LoadBalancer. Handle the HTTP queries of the Azure LB, check the
state of the cluster and provide feedback.

# TODO

- [ ] Make options configurable
- [x] Make it reliable


# Current state

The health check should provide a HTTP response within reasonable
amount of time in the following conditions:

1. Node not running (TCP SYN answered with RST)
2. Node IP ending in black hole (connection timeout)
3. Network trouble (no data, slow data)
4. Node very busy ("slowloris like apache issue")
4. Node not a primary

