
Cli to rewrite `shipyard module blueprints` to automate enabling/disabling of a combination of shipyard modules and their dependencies. 
It also wraps calls to `shipyard` by switching dir/ to the target directory containing the modules definitions

Commands 
```shell
shpMod [cmd] [args]                                                                                                                                                                                                                                            
enable      -- enable modules `shpMod enable [tab] [tab]`
run         -- shipyard run
destroy     -- shipyard destroy
setConfig   -- set $HOME/.shpMod/cfg.hcl `./shpMod setConfig ./config.hcl` or clear `./shpMod setConfig -clear` 
showConfig  -- print $HOME/.shpMod/cfg.hcl `./shpMod showConfig`
```

Optional config.
```hcl
Config {
  // default modules that should be always enabled
  enable = [  // optional
    "network",
  ]
  // skip these for shell completion
  skipShellCompletion = [ // optional
    "database",
    "nats",
    "network",
  ]
}
```