<!--
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-04-07 22:55:18
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-05-13 20:31:10
 * @FilePath: /youtube/readme.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

# SideTube

SideTube is a side project that implements a simple youtube. Everyone can upload videos to the platform and message each video. The system will auto convert uploading mp4 video to webm-dash with different px (480p, 720p, 1080p) and fps(30, 60) for adaptive streaming.

This system has been built with a multi-Golang microserver that is easy to follow Domain Drive Design for Long-term development. On the other hand, using Grpc to connect each server reduces network bandwidth. The encoding part uses FFmpeg with Rabbitmq and Redis distribute lock to encode video mission that is easy to scale.

### Directories

| folder         | introduction                                         |
| :------------- | :--------------------------------------------------- |
| nginx          | Http Router                                          |
| env            | docker-compose env files                             |
| stateful       | config and data of ElasticSearch, Rabbitmq and mongo |
| front          | React18 with [Material-UI](https://mui.com/)         |
| backend        | Golang Api server                                    |
| backend/encode | Encoding video by FFmpeg                             |

### System Architecture

![](https://hackmd.io/_uploads/SysJ7E5Vh.png)

| Video Upload  | Channel Studio   | Video Streaming | Video Message |
| ------------- | :--------------- | --------------- | ------------- |
| upload videos | edit video meta  | provide video   | message       |
|               | video list       |                 | edit message  |
|               | set public state |                 | reply         |
|               |                  |                 | edit reply    |

| Recommend          | User          | JWT generate        | Picture     | Search             |
| :----------------- | :------------ | :------------------ | :---------- | :----------------- |
| home recommend     | login         | create access token | videos plan | video title search |
| relative recommend | logout        |                     |             |                    |
|                    | register      |                     |             |                    |
|                    | get user info |                     |             |                    |
|                    | history       |                     |             |                    |

### Reference

https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example

### Running locally

##### Run cmd

```
mkdir -p ../tempVideo ./stateful/elasticSearch/data  ./stateful/mongo/data ./stateful/rabbitmq/data

docker-compose up
```

```
youtube-mongo-1          | {"t":{"$date":"2023-04-06T04:37:09.312+00:00"},"s":"I",  "c":"STORAGE",  "id":22430,   "ctx":"Checkpointer","msg":"WiredTiger message","attr":{"message":"[1680755829:312543][1:0xffff7f11fd00], WT_SESSION.checkpoint: [WT_VERB_CHECKPOINT_PROGRESS] saving checkpoint snapshot min: 29, snapshot max: 29 snapshot count: 0, oldest timestamp: (0, 0) , meta checkpoint timestamp: (0, 0) base write gen: 1592"}}
youtube-mongo-1          | {"t":{"$date":"2023-04-06T04:38:09.349+00:00"},"s":"I",  "c":"STORAGE",  "id":22430,   "ctx":"Checkpointer","msg":"WiredTiger message","attr":{"message":"[1680755889:349484][1:0xffff7f11fd00], WT_SESSION.checkpoint: [WT_VERB_CHECKPOINT_PROGRESS] saving checkpoint snapshot min: 31, snapshot max: 31 snapshot count: 0, oldest timestamp: (0, 0) , meta checkpoint timestamp: (0, 0) base write gen: 1592"}}
youtube-mongo-1          | {"t":{"$date":"2023-04-06T04:39:14.770+00:00"},"s":"I",  "c":"STORAGE",  "id":22430,   "ctx":"Checkpointer","msg":"WiredTiger message","attr":{"message":"[1680755954:770521][1:0xffff7f11fd00], WT_SESSION.checkpoint: [WT_VERB_CHECKPOINT_PROGRESS] saving checkpoint snapshot min: 33, snapshot max: 33 snapshot count: 0, oldest timestamp: (0, 0) , meta checkpoint timestamp: (0, 0) base write gen: 1592"}}
youtube-elasticsearch-1  | {"@timestamp":"2023-04-06T04:39:14.775Z", "log.level": "WARN", "message":"timer thread slept for [1m/64187ms] on absolute clock which is above the warn threshold of [5000ms]", "ecs.version": "1.2.0","service.name":"ES_ECS","event.dataset":"elasticsearch.server","process.thread.name":"elasticsearch[22065c1285ba][[timer]]","log.logger":"org.elasticsearch.threadpool.ThreadPool","elasticsearch.cluster.uuid":"jQyMreCcRiyO4b9MTgxL1g","elasticsearch.node.id":"jqStWs04Q-qXLWmtWXf-FQ","elasticsearch.node.name":"22065c1285ba","elasticsearch.cluster.name":"docker-cluster"}
youtube-kibana-1         | [2023-04-06T04:39:17.885+00:00][INFO ][status] Kibana is now degraded (was available)
youtube-kibana-1         | [2023-04-06T04:39:35.887+00:00][INFO ][status] Kibana is now available (was degraded)
```

##### docker exex -it < mongo container id > bash

```
mongo --authenticationDatabase admin --username root

use video_meta

db.createUser({
     user: "video_meta",
     pwd: "video_meta",
     roles: [
     "dbOwner",
     ]
 });

```

#### Instell [Migration Cli](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) and run commend

```
brew install golang-migrate

migrate -path "./backend/db_migration/mongo/"  -database 'mongodb://video_meta:video_meta@127.0.0.1:27017/video_meta' up
```

##### Watch it at http://localhost

### Screenshots

![](https://i.imgur.com/7iNWJYq.jpg)

![](https://i.imgur.com/fOOA4K2.jpg)

![](https://i.imgur.com/UWufOgP.png)

### Contributors

[ChatGPT](https://openai.com)
