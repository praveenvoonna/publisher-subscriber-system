# Real-Time Match Data Publisher-Subscriber System

## Glossary

### Term
Description

<font color="red">All terms marked with this color are env configurable</font>

## Objective

The objective of this system is to retrieve ongoing match data from a specified source and transmit the same to subscribed clients in real-time. It utilizes Server-Sent Events (SSE), establishing a unidirectional socket-like connection between the locally running server and the client, allowing data to be pushed from the server and reflected in the client's browser.

## Simplified Steps

1. Fetch match data from the source using an HTTP request and store it in memory/Redis utilizing a cron job at the required interval.
2. Retrieve match data from memory/Redis and send it to the client at the specified interval.

## Assumptions

- Match data source: [https://reqres.in/api/users/2](https://reqres.in/api/users/2)
- Data transmission to clients using SSE via [github.com/antage/eventsource](https://github.com/antage/eventsource)
- Data polling facilitated by [github.com/go-co-op/gocron](https://github.com/go-co-op/gocron)

### Sample Output JSON

Each data stream should deliver data in the following format:

```json
{
    "timestamp": 1683240226,
    "data": {
        "username": "janet.weaver",
        "domain": "reqres.in",
        "name": "janet weaver"
    }
}
```

## Deliverables

The system should provide the following deliverables:

1. **An API Endpoint for Real-Time Data:**
   - The API endpoint should establish a constant stream of data to the client's browser using Server-Sent Events (SSE).
   
2. **Transmission Frequency:**
   - The API should transmit data to the client at regular intervals of every 2 seconds.
   
3. **Data Fetching from Source:**
   - Data retrieval from the source should occur at intervals of every 10 seconds to ensure the most recent match information is obtained.
   
4. **Data Caching with Redis:**
   - Data fetched from the source should be cached in Redis to optimize subsequent data access and minimize source requests.
   
5. **Optional: Logging API Requests:**
   - Optionally, API requests made to the data source can be logged in a database for monitoring, analysis, or auditing purposes.
