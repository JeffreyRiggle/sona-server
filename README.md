The Sona Project
================

Sona (Support Oriented Notification Application) is a open-source project. This project is intended to provide a simple and easy to use web service that allows for the creation and management of incidents.

# Sona

## Overview

Sona is a web service that provides a configurable incident management system. Its original intent was to create a simple and easy to configured pipleline for the maintence of software.
Sona provides

- Options for how incidents are stored.
- Options for how attachments are stored.
- A REST API to manage incidents
- Configuration to allow the invocation of other web services.

## Audience

Sona is recommneded for anyone who wants a free easy to deploy incident tracking system. 


## Getting Started

### Setup

  1. Install go 1.8.3+ (https://golang.org/dl/)
  2. Create a new folder for this project making sure it is in your gopath.
  3. run git clone https://<username>@bitbucket.org/JeffreyRiggle/sona.git
  
### Building
  1. Run go get (make sure to grab dependencies)
  2. Run go build in source folder
  
### Running

  1. Run the src(.exe) file with elevated privledges.
  
## Configuration

Sona offers a couple of configuration options for the application at runtime. In order to run with a configuration run src with the file (inculding path) as a command argument.

Some configuration options to consider would be the incidentmanagertype and filemanagertype. For more information on these configurations see configuration.go
