# GraphQL 查詢示例

## GraphQL Playground

啟動服務器後，訪問 `http://localhost:8080/graphql` 可以使用 GraphQL Playground 進行交互式查詢。

## 查詢示例

### 1. 獲取所有會員

```graphql
query {
    users {
        id
        name
        email
    }
}
```

### 2. 根據 ID 獲取單個會員

```graphql
query {
    user(id: 1) {
        id
        name
        email
    }
}
```

## Mutation 示例

### 3. 創建新會員

```graphql
mutation {
    createUser(name: "張三", email: "zhangsan@example.com") {
        id
        name
        email
    }
}
```

### 4. 更新會員信息

```graphql
mutation {
    updateUser(id: 1, name: "李四", email: "lisi@example.com") {
        id
        name
        email
    }
}
```

### 5. 部分更新會員信息（只更新名稱）

```graphql
mutation {
    updateUser(id: 1, name: "王五") {
        id
        name
        email
    }
}
```

### 6. 刪除會員

```graphql
mutation {
    deleteUser(id: 1)
}
```

## 完整示例

### 創建會員並查詢

```graphql
mutation {
    createUser(name: "測試用戶", email: "test@example.com") {
        id
        name
        email
    }
}

query {
    users {
        id
        name
        email
    }
}
```
