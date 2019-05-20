# go-soccer

## RESTful 
| Resource  | GET (read)        | POST (create)       | PULL (update)       | DELETE              |
|-----------|-------------------|---------------------|---------------------|---------------------|
| /teams    | Not allowed (405) | Create a new team   |  Not allowed (405)  |  Not allowed (405)  |
| /teams/1  |  Returns team 1   |  Not allowed (405)  |  Not allowed (405)  |  Not allowed (405)  |

# How to start

- initialize db

    ```
    go run main.go db init
    ```

- create a staker account

    ```
    go run main.go staker new
    ```

    this command will generate a keystore and suggest an onion hash. Add those values in the following config.yaml section

    ```
    stakers:
    	- address: <keystore_address>
    	  onion:   <onion_hash>
    ```

repeat this process if you want to generate more accounts

- enroll as a staker
    ```
    go run main.go staker enroll <keystore_address>
    ```

- run the service
    ```
    go run main.go service start
    ```

