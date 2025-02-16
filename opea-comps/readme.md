## Running Ollama third party service


```
NO_PROXY=localhost LLM_ENDPOINT_PORT=8008 LLM_MODEL_ID="llama3.2" host_ip=192.168.0.226 docker compose up
```

### Ollama API

Once the Ollama sv is running we can make API calls to the API

## Generate a request

curl http://localhost:8008/api/pull -d '{"model": "llama3.2:1B"}'

curl http://localhost:8008/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?"
}'