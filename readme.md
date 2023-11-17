# Tender Matching System

## 系統架構圖
```mermaid

graph LR


subgraph tinder matching system
    subgraph gin
        AddSinglePersonAndMatch
        RemoveSinglePerson
        QuerySinglePeople
    end
    subgraph match_system
        maleHeightIndex
        faleHeightIndex
        nameIndex
        singlePersons
    end
    gin --> match_system

end


```
## 說明

主要透過 Gin 對外開放三隻 API 分別會是 AddSinglePersonAndMatch、RemoveSinglePerson 與 QuerySinglePeople，這三隻 API 收到 request 後會將參數整理後交由 match_system 進行後續處理。

match_system 本身將 user 資訊儲存與 singlePersons 在儲存與刪除 user 資訊時會同時異動 maleHeightIndex、faleHeightIndex 與 nameIndex。

其 match_system 相關作業皆都依賴於 maleHeightIndex、faleHeightIndex 與 nameIndex。

## API 時序圖與時間複雜度

### AddSinglePersonAndMatch

```mermaid
graph TD
A[start] --> B{ValidatePerson}
B -->|參數錯誤| e[return]
B -->|參數正確| C{確認person是否已經存在}
C -->|已存在| e[return]
C -->|不存在| D[將 person 存入 singlePersons]
D --> E[將 ID 根據 gender 與 height 塞入相對應的 index]
E --> e[return]

```


過程主要都是對 map 進行操作，golang map 操作平均時間複雜度為 O(1)，故整體時間複雜度為 O(1)

---

### RemoveSinglePerson

```mermaid
graph TD
A[start] --> B{參數檢查}
B -->|參數錯誤| e[return]
B -->|參數正確| C{確認person是否已經存在}
C -->|不存在| e[return]
C -->|已存在| D[根據 gender 與 height 刪除相對應的 index]
D --> E[根據 ID 從 singlePersons 刪除]
E --> e[return]

```

DeleteHeightIndex 的內容會針對 height 尋找相對應的 slice ，並將 id 從該 slice 刪除，這裡的時間複雜度會為 O(k)，k 為 slice 的長度，除此之外 RemoveSinglePerson 其餘操作皆為 O(1)，因此時間複雜度為 O(k)。

---

### QuerySinglePerson

```mermaid
graph TD
A[start] --> B{參數檢查}
B -->|參數錯誤| e[return]
B -->|參數正確| C[根據 query 參數找尋可能的 index ]
C --> E[根據 index 篩選出條件內的適合人選]
E --> e[return]

```

QuerySinglePerson 的內容會針對 query 的條件尋找可能符合的 index ，再透過這些 index 做細部的篩選，FindHeightIdx 時間複雜度為 O(N * K) 其中 N 是 heightIndex 中的鍵值對數量，K 是每個切片的平均長度。

最後篩選出條件內的適合人選得部分會是 FindHeightIdx 的結果所以也是 O(N * K)，所以整個 API 時間複雜度為 O(N * K) 

---

# Generate API Doc
```
make doc
```

# Docker
```
make docker
docker run -p 8080:80 -d tinder_match_system  
```


