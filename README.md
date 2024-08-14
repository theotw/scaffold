# Scaffold Project
## Overview
The purpose of this project is to keep a collection of code that I keep remaking in every project

Specifically
- [ ] Logging
- [ ] Configuration
- [ ] Error Handling
- [ ] And some misc functions like GetEnv with a default value


## Logging
Is a simple pluggable logger that can use logrus or zap

# Errorutils
Is a simple package that has a wrapper function for returns and common errors prefefined

# Config
I wanted a standard config package that can read from a file or env variables or AWS Secrets
But I dont like just maps of config so it will read the config into a struct defined by the app
Some conventions are used
* Config is all strings
* Config in the env vars are prefixed with a user selected prefix. like MYAPP_ or something
* The structure config items need a json tag to map to the env var name minus the prefix and lower case

There are some common env vars read
* AWS_DEFAULT_REGION
* AWS_ACCESS_KEY_ID
* AWS_SECRET_ACCESS_KEY
* AWS_PROFILE
  