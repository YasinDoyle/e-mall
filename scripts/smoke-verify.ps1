param(
    [string]$BaseUrl = "http://127.0.0.1:5001/api/v1",
    [Parameter(Mandatory = $true)][string]$UserName,
    [Parameter(Mandatory = $true)][string]$Password,
    [Parameter(Mandatory = $true)][string]$PayKey,
    [Parameter(Mandatory = $true)][int]$ProductId,
    [Parameter(Mandatory = $true)][int]$BossId,
    [Parameter(Mandatory = $true)][int]$AddressId,
    [Parameter(Mandatory = $true)][decimal]$Money,
    [int]$Num = 1
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers,
        [object]$Body
    )

    $jsonBody = $null
    if ($null -ne $Body) {
        $jsonBody = $Body | ConvertTo-Json -Depth 6
    }

    return Invoke-RestMethod -Method $Method -Uri $Url -Headers $Headers -ContentType "application/json" -Body $jsonBody
}

Write-Host "[1/4] 登录获取 token..."
$loginResp = Invoke-Api -Method Post -Url "$BaseUrl/user/login" -Headers @{} -Body @{
    user_name = $UserName
    password  = $Password
}

$accessToken = $loginResp.data.access_token
$refreshToken = $loginResp.data.refresh_token
if (-not $accessToken -or -not $refreshToken) {
    throw "登录失败，未拿到 access_token/refresh_token。"
}

$headers = @{
    access_token  = $accessToken
    refresh_token = $refreshToken
}

Write-Host "[2/4] 创建订单..."
$createResp = Invoke-Api -Method Post -Url "$BaseUrl/orders/create" -Headers $headers -Body @{
    product_id = $ProductId
    num        = $Num
    address_id = $AddressId
    money      = $Money
    boss_id    = $BossId
}

if (-not $createResp.status -or $createResp.status -ne 200) {
    throw "创建订单失败: $($createResp | ConvertTo-Json -Depth 6)"
}

Write-Host "[3/4] 查询最新订单列表，拿订单 ID..."
$listResp = Invoke-Api -Method Get -Url "$BaseUrl/orders/list?page_num=1&page_size=10" -Headers $headers -Body $null
$orderId = $listResp.data.item[0].id
if (-not $orderId) {
    throw "未从订单列表里找到订单 ID。"
}

Write-Host "[4/4] 连续支付两次，验证幂等..."
$firstPay = Invoke-Api -Method Post -Url "$BaseUrl/paydown" -Headers $headers -Body @{
    order_id   = $orderId
    product_id = $ProductId
    boss_id    = $BossId
    num        = $Num
    key        = $PayKey
}

$secondPay = Invoke-Api -Method Post -Url "$BaseUrl/paydown" -Headers $headers -Body @{
    order_id   = $orderId
    product_id = $ProductId
    boss_id    = $BossId
    num        = $Num
    key        = $PayKey
}

Write-Host "第一次支付响应:"
$firstPay | ConvertTo-Json -Depth 6
Write-Host "第二次支付响应:"
$secondPay | ConvertTo-Json -Depth 6

Write-Host "烟测结束。重点看第二次支付是否被拒绝。"