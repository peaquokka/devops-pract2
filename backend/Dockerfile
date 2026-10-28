# ---------- Этап 1: сборка ----------
FROM golang:1.22-alpine AS builder

# go mod download использует GOMODCACHE, чтобы кэшировать модули
WORKDIR /build

# Сначала зависимости — слой кэшируется, пока go.mod не изменится.
# Это ускоряет пересборку: при правке кода зависимости не качаются заново.
# go.sum генерируется внутри контейнера через go mod tidy — поэтому копируем только go.mod.
COPY go.mod ./
RUN go mod download && go mod tidy

# Потом исходники — этот слой перезаписывается при любом изменении кода.
COPY . .
RUN go build -mod=mod -ldflags="-s -w" -o /out/app ./main.go

# ---------- Этап 2: исполнение ----------
# distroless: без shell, nonroot
FROM gcr.io/distroless/static-debian12:nonroot

# nonroot-вариант distroless запускается под пользователем nobody/nogroup
WORKDIR /app
COPY --from=builder /out/app /app/app

EXPOSE 8080

# ENTRYPOINT — бинарник, CMD — аргументы по умолчанию
ENTRYPOINT ["/app/app"]