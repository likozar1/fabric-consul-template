## Generating Fabric data file from consul catalog using consul-template

Get data from consul catalog using GET request. Process data and generate output python file.

## How it works
 
Consul_template will run fabric_consul binary when something change in consul catalog. 

### Usage

import generated file inside fabric file

```python
	from [output-file] import *
```