module ds-uoe-vash/prime-app

go 1.25.3

// install the dep from our local folder instead of trying to fetch it online
replace ds-uoe-vash/dfs => ../dfs

require ds-uoe-vash/dfs v0.0.0
