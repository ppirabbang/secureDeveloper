# Go Secure Coding Practice

보안 코딩 연습을 위한 시작 프로젝트입니다.

처음부터 구조를 예쁘게 나누기보다, `cmd/server/main.go` 하나에 코드를 모아 둔 상태에서
먼저 흐름을 이해하고 직접 분리 기준을 고민할 수 있게 만드는 것이 목적입니다.

지난 과제와 수업에서 설명했던 내용, 그리고 전달된 가이드 문서를 떠올리면서
어떤 기능부터 구현하고 어떤 기준으로 나눌지 스스로 판단해 보세요.

## 프로젝트 목적

- 로그인 흐름을 먼저 이해하기
- 더미 API를 실제 동작 코드로 바꿔 보기
- 게시판과 금융 기능을 단계적으로 채워 보기
- 입력 검증, 인증 확인, 권한 검사, 응답 설계를 직접 고민해 보기
- 구현 후 어떤 코드들을 묶어 리팩터링할지 판단해 보기

## 현재 상태

- 서버 코드는 현재 `cmd/server/main.go` 하나에 들어 있습니다.
- 전체적인 기능이 구현되어 있습니다
- 추후 관리를 편의를 위해 코드 모듈화를 진행했습니다
  - 기본적으로 middleware - handler - service - ext(db, cache, logger) 로 흐름을 추적하면 됩니다
- traceId 기반으로 로그를 추적하도록 설정되어있습니다
- 
```aiignore
ex) 쉘 명령으로 확인 
cat logs/app.log | grep { 특정 요청의 traceId } 
```

## 주요 API

인증
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /api/auth/withdraw`

사용자
- `GET /api/me`

게시판
- `GET /api/posts`
- `GET /api/posts/:id`
- `POST /api/posts`
- `PUT /api/posts/:id`
- `DELETE /api/posts/:id`

금융
- `POST /api/banking/deposit`
- `POST /api/banking/withdraw`
- `POST /api/banking/transfer`

## 참고 파일
- `schema.sql`
- `seed.sql`
- `query_examples.sql`

## 실행 방법

프로젝트 루트에서 실행합니다.

```powershell
go run ./cmd/server
```
처음 상태로 다시 시작하고 싶으면 `app.db`를 지운 뒤 다시 실행하면 됩니다.

