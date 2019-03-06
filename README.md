# go-soccer

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

