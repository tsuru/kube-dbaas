# kube-dbaas
Experimental (Hackday) project to provide a tsuru service that maintains database operators


# Inspiration

We are inspirated in our main project https://github.com/globocom/database-as-a-service that provides thousands of database vms inside Globo

# Motivation

Containers are cheaper and more scalable way to deploy software in nowadays, this proof of concept will use battle tested kubernetes operators to maintain a simple API.

# Supported Operators
- mongodb/mongodb-kubernetes-operator
- spotahome/redis-operator
