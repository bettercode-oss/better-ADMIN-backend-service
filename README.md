# better ADMIN Back-end Service

Golang 으로 구현한 better ADMIN Back-end Service

## 도커

### 도커 이미지 빌드
```
docker build --no-cache --rm=true --tag bettercode2016/better-admin-backend-service:latest .
```

### 도커 Hub 업로드
```
docker push bettercode2016/better-admin-backend-service:latest
```

### 실행 
```
docker run -d -p 2016:2016 bettercode2016/better-admin-backend-service
```
