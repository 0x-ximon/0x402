# 0x402
0x402 aims to simplify x402 Integration on Aptos for Go web applications. It provides a simple and easy-to-use CLI for kick starting a new Go x402 application as well as middleware that you can plug and play into your existing Go. 

## As a Dependency

To use 0x402 as a dependency in your Go project, run the following command:

```bash
go get -u github.com/0x-ximon/0x402
```

### Configuration and Usage

You can configure 0x402 by setting environment variables or by using a configuration file. Here's an example of how to configure 0x402 using environment variables:

```go 
import (
		"github.com/0x-ximon/0x402/middlewares"
)

func main() {
		mux := http.NewServeMux()
		addr := "http://localhost:8080"
	
		cfg := middlewares.PaymentConfig{
				Amount:      10000, // 0.01 USDC
				Receiver:    "YOUR_ADDRESS",
				Description: "Aptos is Awesome",
				Asset:       "0x69091fbab5f7d635ee7ac5098cf0c1efbe31d68fec0f2cd565e8d168daf52832", // USDC on Aptos
		}
	
		// The Guard accepts a configuration struct and provides 2 Paywall middleware
		// functions that are compatible with the Go standard http package. You can
		// set it to nil and use the default configuration
		guard := middlewares.NewGuard(nil)
	
		// Helper function to create chain of middlewares
		chain := middlewares.NewChain(
				// The StandardPaywall middleware is used to protect the entire
				// application and takes the payment config
				guard.StandardPaywall(cfg),
		)
	
		server := http.Server{
				Addr:    addr,
				Handler: chain(mux),
		}
	
		log.Printf("Starting server on %s", addr)
		server.ListenAndServe()
}
```

> [!WARNING]
> Only the `StandardPaywall` middleware has been fully implemented. The `ResourcePaywall` middleware is currently in development and will be available soon.

## Installation

Although this feature is currently in development, you can install this package as a Binary and use it to protect your application. To install 0x402, run the following command:

```bash
go install github.com/0x-ximon/0x402@latest
```
<img width="996" height="502" alt="image" src="https://github.com/user-attachments/assets/965acaef-157b-45c2-99a9-329950ae6cf2" />

### Commands

* `0x402 help`: Display help information for the `0x402` command.
* `0x402 init`: Initialize a new Go project with x402 configured on Aptos.
* `0x402 add`: Generates and adds a new Aptos resource to your project.
* `0x402 inspect`: Inspects an Aptos resource and displays its details.
* `0x402 query`: Queries the Aptos blockchain for a resource.
