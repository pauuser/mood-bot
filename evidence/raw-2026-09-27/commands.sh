# Команды, которыми получены отчёты в этой папке.
# Запуск из корня mood-bot на коммите 81ed500. Нужен только Docker.

# gitleaks.txt
docker run --rm -v "$PWD":/repo zricethezav/gitleaks:v8.30.1 git /repo -v

# semgrep.txt
docker run --rm -v "$PWD":/src -w /src semgrep/semgrep:1.177.0 \
  semgrep scan --config p/golang --config p/dockerfile --config p/docker-compose --metrics=off .

# govet-govulncheck.txt
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c \
  'go vet ./...; go install golang.org/x/vuln/cmd/govulncheck@v1.1.3 && govulncheck -show verbose ./...'

# govulncheck-go1.25.14.txt
docker run --rm -e GOTOOLCHAIN=local -v "$PWD":/src -w /src golang:1.25.14 sh -c \
  'go install golang.org/x/vuln/cmd/govulncheck@v1.1.4 && govulncheck ./...'

# probe-tests.txt (probe_test.go.txt копируется в internal/usecases/usecases_impl/probe_test.go)
docker run --rm -v "$PWD":/src -w /src golang:1.21 \
  go test -v -count=1 ./internal/usecases/usecases_impl/

# compose-run.txt
TELEGRAM_BOT_TOKEN=123456:FAKE docker compose up --build --abort-on-container-exit

# compose-chown-linux.txt (каталог данных принадлежит root, как у bind mount на Linux)
docker run --rm --mount type=tmpfs,destination=/app/data mood-bot-mood-bot sh -c \
  'mkdir -p /app/data && chown -R appuser:appuser /app/data && ./main'

# git-history.txt
git log --format='%h %ad author=%an committer=%cn %s' --date=short
curl -s 'https://api.github.com/repos/pauuser/mood-bot/pulls?state=all' | jq length
curl -s 'https://api.github.com/repos/pauuser/mood-bot/actions/runs' | jq .total_count
