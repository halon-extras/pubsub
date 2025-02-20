# Google Cloud Pub/Sub plugin

This plugin is a wrapper for the [Google Cloud Pub/Sub](https://pkg.go.dev/cloud.google.com/go/pubsub) library.

## Installation

Follow the [instructions](https://docs.halon.io/manual/comp_install.html#installation) in our manual to add our package repository and then run the below command.

### Ubuntu

```
apt-get install halon-extras-pubsub
```

### RHEL

```
yum install halon-extras-pubsub
```

## Configuration
For the configuration schema, see [pubsub.schema.json](pubsub.schema.json). Below is a sample configuration.

Note that the [JSON key file](https://cloud.google.com/iam/docs/managing-service-account-keys) setting (`credentials`) is optional and [Google Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials) will be used if omitted.

### smtpd.yaml

```
plugins:
  - id: pubsub
    config:
      projects:
        - id: project-id
          topics:
            - topic1
          # credentials: /path/to/keyfile.json
```

## Exported functions

These functions needs to be [imported](https://docs.halon.io/hsl/structures.html#import) from the `extras://pubsub` module path.

### pubsub_publish(project, topic, message)

Publish a message to a topic asynchronously (Supports HSL suspend functionality).

**Params**

- project `string` - The ID of the project as configured in `smtpd.yaml`
- topic `string` - The topic as configured in `smtpd.yaml`
- message `string` - The message to publish

**Returns**

A successful publish will return an associative array with a `id` key that contains the ID of the message. On error an additional `error` key will be provided with the error message.

**Example**

```
import { pubsub_publish } from "extras://pubsub";
echo pubsub_publish("project-id", "topic1", "This is a sample message.");
```
