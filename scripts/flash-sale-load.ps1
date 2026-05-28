param(
    [string]$BaseUrl = "http://127.0.0.1:5001/api/v1",
    [Parameter(Mandatory = $true)][string]$AccessToken,
    [Parameter(Mandatory = $true)][string]$RefreshToken,
    [Parameter(Mandatory = $true)][int]$ProductId,
    [Parameter(Mandatory = $true)][int]$AddressId,
    [Parameter(Mandatory = $true)][string]$PayKey,
    [int]$Requests = 20
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$jobScript = {
    param($Url, $AccessTokenValue, $RefreshTokenValue, $Product, $Address, $KeyValue)

    $headers = @{
        access_token  = $AccessTokenValue
        refresh_token = $RefreshTokenValue
    }

    try {
        $resp = Invoke-RestMethod -Method Post -Uri "$Url/flash_sale/skill" -Headers $headers -ContentType "application/json" -Body (@{
            product_id = $Product
            address_id = $Address
            key        = $KeyValue
        } | ConvertTo-Json)

        [PSCustomObject]@{
            ok   = $true
            body = $resp
        }
    }
    catch {
        [PSCustomObject]@{
            ok   = $false
            body = $_.Exception.Message
        }
    }
}

Write-Host "启动 $Requests 个并发秒杀请求..."
$jobs = @()
for ($i = 0; $i -lt $Requests; $i++) {
    $jobs += Start-Job -ScriptBlock $jobScript -ArgumentList $BaseUrl, $AccessToken, $RefreshToken, $ProductId, $AddressId, $PayKey
}

$jobs | Wait-Job | Out-Null
$results = $jobs | Receive-Job
$jobs | Remove-Job | Out-Null

$success = @($results | Where-Object { $_.ok -and $_.body.status -eq 200 }).Count
$failed = @($results | Where-Object { -not $_.ok -or $_.body.status -ne 200 }).Count

Write-Host "成功请求数: $success"
Write-Host "失败请求数: $failed"
Write-Host "样例响应:"
$results | Select-Object -First 5 | ConvertTo-Json -Depth 8