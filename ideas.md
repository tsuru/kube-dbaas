# Ideas

Some random thoughts on how we could evolve this API to be more flexible, without the need to write custom code for each database we want to offer.

## Kube DBaaS Configuration API

POST, GET /api/v1/database
	- name: string
	- templates: []string
	- fields: [
		- Name:
		- Type: []
		- Default:
	]

GET, PUT /api/v1/database/:name
	- name: string
	- templates: []string
	- fields: [
		- Name:
		- Type: []
		- Default:
	]

	==> Creates a CRD with the values passed, this configures a new database type which will be listed on `/resources/plans` call.
	- templates is a collection of helm chart like strings.
	- fields is the template schema, ie. all field that can be filled when creating an instance.

/api/v1/plans
	- plan name
	- resource limits
	example: c1m2 => cpu: 1000m memory: 2Gi

	==> Creates a CRD with resources limits, `/resource/plans` does the cartesian product between database types and plans.

## Tsuru Service Instance API

/resources/plans
/resources
/resources/:instance
/resources/:instance/autoscale
/resources/:instance/bind
/resources/:instance/bind-app
/resources/:instance/info
/resources/:instance/plans
/resources/:instance/status