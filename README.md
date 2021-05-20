# better ADMIN Back-end Service

Golang 으로 구현한 better ADMIN Back-end Service


## 데이터베이스

### Sqlite
별도 설정을 하지 않는다면 기본적으로는 Sqlite file 데이터베이스를 사용한다.

### MySQL
* 데이터 베이스 생성
```sql
-- 아래 데이터베이스명은 예시
CREATE SCHEMA IF NOT EXISTS `better_admin` DEFAULT CHARACTER SET utf8mb4;
```

* 애플리케이션 실행 환경 변수 설정
```
BETTER_ADMIN_DB_HOST=localhost:3306
BETTER_ADMIN_DB_DRIVER=mysql
BETTER_ADMIN_DB_NAME=better_admin
BETTER_ADMIN_DB_USER=root
BETTER_ADMIN_DB_PASSWORD=1111
```

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
