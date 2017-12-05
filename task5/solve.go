package main

import(
    "encoding/json"
    "net/http"
    "fmt"
    "sort"
    "strconv"

    "github.com/gorilla/mux"
)
type messageHandler struct{
    message string
}
type myItem struct{
  key string `json:"key"`
}

var (
  itemsStore  = make(map[int]myItem)
  key = 0
)

func GetItems(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  items := make([]myItem, 0, len(itemsStore))

  var keys []int
  for k := range itemsStore {
    keys = append(keys, k)
  }
  sort.Ints(keys)

  for _, k := range keys {
    items = append(items, itemsStore[k])
  }

  j, err := json.Marshal(items)
  if err != nil {
    panic(err)
  }
  w.WriteHeader(http.StatusOK)
  w.Write(j)
}

func UpdateItems(w http.ResponseWriter, r *http.Request){
  vars := mux.Vars(r)
  k, err := strconv.Atoi(vars["key"])
  if err != nil {
    panic(err)
  }
  var itemToUpdate myItem
  err = json.NewDecoder(r.Body).Decode(&itemToUpdate)
  if err != nil{
    panic(err)
  }
  if item, ok := itemsStore[k]; ok {
    itemToUpdate.key = item.key
    delete(itemsStore, k)
    itemsStore[k] = itemToUpdate
  }
  w.WriteHeader(http.StatusNoContent)
}

func (m *messageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, m.message)
}

func main(){
    ///mux := http.NewServeMux()
    router := mux.NewRouter()
    apiRouter := router.PathPrefix("solve").Subrouter()
    apiRouter.HandleFunc("/",GetItems).Methods("GET")
    apiRouter.HandleFunc("/{key}", UpdateItems).Methods("UPDATE")
    //apiRouter.HandleFunc("/",AddTodoItems).Methods("POST")
    //fs := http.FileServer(http.Dir("solve"))
    //mux.Handle("/", fs)
    //mh1 := &messageHandler{"adsasd"}
    //mux.Handle("/welcome", mh1)
    //mux.HandleFunc("/", )
    http.ListenAndServe(":8082", router)
}
