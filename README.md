
The `aks-logs` tool will help you to query the Container logs running inside the Azure Kubernetes Service cluster.

## Getting Started

### Prerequisite

1. [Azure CLI installed](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)

### Installation

#### Linux

1. Download

    ```bash
    curl -LO https://github.com/Mirage20/aks-logs/releases/latest/download/aks-logs-linux-x64.tar.gz
    ```
2. Extract

    ```bash
    tar -xzvf aks-logs-linux-x64.tar.gz
    ```
3. Install

    ```bash
     sudo mv ./aks-logs /usr/local/bin/aks-logs
    ```

#### MacOS

1. Download

    ```bash
    curl -LO https://github.com/Mirage20/aks-logs/releases/latest/download/aks-logs-darwin-x64.tar.gz
    ```
2. Extract

    ```bash
    tar -xzvf aks-logs-darwin-x64.tar.gz
    ```
3. Install

    ```bash
     sudo mv ./aks-logs /usr/local/bin/aks-logs
    ```

### Usage

1. [Login to your Azure account using the CLI](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli)
   
2. Execute `aks-logs` command with required flags

    ```bash
    # Check logs for my-service workload
    aks-logs -subscriptionId=00000000-0000-0000-0000-000000000000 \
       -resourceGroupName=my-res-group -clusterName=cluster-1 \
       -namespace=default -workloadName=my-service -containerName=my-service
    ```

3. Run `aks-logs --help` for more information

    ```text
    Usage of aks-logs:
      -clusterName string
            AKS Cluster name. Required
      -containerName string
            Container name of the Kubernetes workload. Required
      -contains string
              Filter log entries by list of contents if specified. Example: -contains=val1,val2
      -endTime string
            End time of the log entries. Default to (current time) if not specified. Format: 2006-01-02T15:04:05.999999999+07:00
      -maxRecords int
            Maximum number of log entries to query. Default to 1000
      -namespace string
            Kubernetes namespace of the workload. Required
      -resourceGroupName string
            Resource group name of the AKS Cluster. Required
      -showDescending
            Output the logs in descending order based on generated time
      -showQuery
            Output the Kusto query
      -showRaw
            Output the raw results
      -startTime string
            Start time of the log entries. Default to (current time - 24h) if not specified. Format: 2006-01-02T15:04:05.999999999+07:00
      -subscriptionId string
            Subscription id of the AKS Cluster. Required
      -version
            Output version information
      -workloadName string
            Kubernetes workload name. Required
    ```
