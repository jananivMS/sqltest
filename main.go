package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jananiv/sqltest/config"
	"github.com/jananiv/sqltest/sql"
)

var resourcegroup = "Janani-testRG" //Change this to your resource group that the service principal has access to
var location = "westus"
var servername = "sqlsrvrand124"
var adminuser = "iamadmin"
var adminpassword = "$tress12345"

func main() {
	config.ParseEnvironment()

	//Creates SQL Server Async. Gets Created but the func never ends and says "async op still running"
	//CreateServerAsyncWithoutGet()

	//Creates SQL server Async and checks using GetServer for the status. Using the GetServer does help break the infinite async loop above
	CreateServerAsyncWithGet()

	//Update the location --> FAIL: says resource with same name cannot be created in new location
	UpdateLocation("eastus")

	//Update resourcegroup --> FAIL: Says resource with same name exists and so need to use a new name
	UpdateResourceGroup("Janani-TestRG2")

	//Update admin user --> FAIL: Says it updated successfully but doesnt seem to change the adminuser when i look at portal
	UpdateAdminUser("iamadmin2")
}

func CreateServerAsyncWithoutGet() {
	ctx := context.Background()

	for {
		_, err := sql.CreateServer(ctx, resourcegroup, location, servername, adminuser, adminpassword)
		//_,err := acr.DeleteRegistry(ctx, "Janani-testRG", "JananivAcr123")
		if err != nil {
			if strings.Contains(err.Error(), "asynchronous operation has not completed") {
				fmt.Println("Async op ongoing...")
				continue
			}
			fmt.Println("error is %s", err)
			break
		}
	}

}

func CreateServerAsyncWithGet() {
	ctx := context.Background()

	for {
		server, err := sql.GetServer(ctx, resourcegroup, servername)
		if err != nil || *server.State != "Ready" {
			_, err := sql.CreateServer(ctx, resourcegroup, location, servername, adminuser, adminpassword)
			//_,err := acr.DeleteRegistry(ctx, "Janani-testRG", "JananivAcr123")
			if err != nil {
				if strings.Contains(err.Error(), "asynchronous operation has not completed") {
					fmt.Println("Async op ongoing...")
					continue
				}
				fmt.Println("error is %s", err)
				break
			}
		}
		fmt.Println("Server provisioned successfully")
		break
	}
}

func UpdateLocation(newlocation string) {
	ctx := context.Background()

	_, err := sql.CreateServer(ctx, resourcegroup, newlocation, servername, adminuser, adminpassword)
	if err != nil {
		fmt.Println("error is ", err.Error())
	}
}

func UpdateResourceGroup(newresourcegroup string) {
	ctx := context.Background()

	for {
		server, err := sql.GetServer(ctx, newresourcegroup, servername)
		if err != nil || *server.State != "Ready" {
			_, err := sql.CreateServer(ctx, newresourcegroup, location, servername, adminuser, adminpassword)
			if err != nil {
				if strings.Contains(err.Error(), "asynchronous operation has not completed") {
					fmt.Println("Async op ongoing...")
					continue
				}
				fmt.Println("error is", err.Error())
				break
			}
		}
		fmt.Println("Server updated successfully")
		break
	}
}

func UpdateAdminUser(newadminuser string) {
	ctx := context.Background()

	for {
		server, err := sql.GetServer(ctx, resourcegroup, servername)
		if err != nil || *server.State != "Ready" {
			_, err := sql.CreateServer(ctx, resourcegroup, location, servername, newadminuser, adminpassword)
			if err != nil {
				if strings.Contains(err.Error(), "asynchronous operation has not completed") {
					fmt.Println("Async op ongoing...")
					continue
				}
				fmt.Println("error is", err.Error())
				break
			}
		}
		fmt.Println("Server updated successfully")
		break
	}
}
