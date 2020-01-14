package config

import "time"

const ReddisPort = 6379

const MainAppPort = 8080
const NotifsAppPort = 8081
const ChatAppPort = 8090

const AuthServicePort = 8082
const RecommendsServicePort = 8083
const NotifsServicePort = 8084

const ChatWorkers = 3

// const Hostname = "localhost"

const Hostname = "postgres"

// const Hostname = "postgresql://postgres:newPassword@clair_postgres:5432?sslmode=disable"
const Database = "hh"
const User = "postgres"
const Password = "newPassword"

// Time allowed to write a message to the peer.
const WriteWait = 10 * time.Second

// Time allowed to read the next pong message from the peer.
const PongWait = 50 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
const PingPeriod = (PongWait * 9) / 10
const MaxMessageSize = 512

const PublicDir = "/static"

// const MAXUPLOADSIZE = 5 * 1024 * 1024 // 1 mb
