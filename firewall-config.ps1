# 配置Windows防火墙规则

# 启用Windows防火墙
Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled True

# 允许HTTP和HTTPS流量
New-NetFirewallRule -DisplayName "Allow-HTTP" -Direction Inbound -LocalPort 80 -Protocol TCP -Action Allow
New-NetFirewallRule -DisplayName "Allow-HTTPS" -Direction Inbound -LocalPort 443 -Protocol TCP -Action Allow

# 允许Grafana监控服务
New-NetFirewallRule -DisplayName "Allow-Grafana" -Direction Inbound -LocalPort 3000 -Protocol TCP -Action Allow

# 允许后端API服务
New-NetFirewallRule -DisplayName "Allow-Backend-API" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow

# 允许Docker相关端口
New-NetFirewallRule -DisplayName "Allow-Docker" -Direction Inbound -LocalPort 2375,2376 -Protocol TCP -Action Allow

# 允许DNS查询
New-NetFirewallRule -DisplayName "Allow-DNS" -Direction Outbound -RemotePort 53 -Protocol UDP -Action Allow

# 配置Docker网络访问
New-NetFirewallRule -DisplayName "Allow-Docker-Network" -Direction Inbound -Program "C:\Program Files\Docker\dockerd.exe" -Action Allow

# 配置防火墙日志记录
Set-NetFirewallProfile -Profile Domain,Public,Private -LogAllowed True -LogBlocked True -LogIgnored True
Set-NetFirewallProfile -Profile Domain,Public,Private -LogFileName "%systemroot%\system32\LogFiles\Firewall\pfirewall.log"

# 阻止不必要的入站连接
Set-NetFirewallProfile -DefaultInboundAction Block -Profile Public

# 允许已建立的连接
Set-NetFirewallProfile -AllowInboundRules True -AllowLocalFirewallRules True -Profile Domain,Private

# 输出配置结果
Write-Host "防火墙规则配置完成"
Get-NetFirewallProfile | Format-Table Name,Enabled,DefaultInboundAction,DefaultOutboundAction