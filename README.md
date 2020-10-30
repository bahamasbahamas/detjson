# detjson 
## Go-Lib for marshalling alphabetically and deterministically 


Output:
    
    {
      "aa":..., (int)
      "bb":..., (string)
      "vv":...,
      ...,
      "xxxx":[
          "aaa":{...},  (map)
          "bbb":{
              "aa":..., (string)
              "bbb":..., (int)
              "ccc":..., (int)
              ...
              "aaa": {...}, (map)
              "bbb":[...],  (array)
              ...
            },  (map)
            ...
          },
          "ccc":{...}, (map)
          "ddd":{...},  (map)
          "eee":{...},  (map)
          ...
      ]
    }

#### Requirements

go-version from app/go.mod has been installed 

#### How to use it (stores the output in the file test.json)
    marshaller := detjson.NewMarshaller("JSONString")
   	err := marshaller.UnMarshal()
   	if err != nil {
   		log.Fatal(err)
   	}
   	err = marshaller.MarshalOrdered()
   	if err != nil {
   		log.Fatal(err)
   	}
   	// Erstelle Datei für den Output.
   	f, err := os.Create("test.json")
   	if err != nil {
   		log.Fatal(err)
   		return
   	}
   	defer func() {
   		err := f.Close()
   		if err != nil {
   			log.Fatal(err)
   		}
   	}()
   	_, err = f.WriteString(marshaller.GetJSONString())
   	if err != nil {
   		fmt.Println(err)
   		return
   	}
    
