# collins2dynect

Tool to update dynECT hosted domains from collins assets


collins2dynect will:

1. **delete all A records from the given domain**
2. Request all assets from collins
3. Add a DNS record for each address a asset has
4. Add a DNS record [alias].[domain] for each alias in `DNS_ALIASES`
5. Publish the zone

## prerequisites
1. collins-shell
  * gem install collins_shell

## Aliases
To add a record 'foo.[domain]' pointing to asset 'bar's address from
pool 'dmz', set the `DNS_ALIASES` attribute like this:

    collins-shell asset set_attribute DNS_ALIASES "foo@dmz" --tag="bar"

DNS_ALIASES is a whitespace separated list.


## Hostname format

The format (hardcoded for now):

    [PRIMARY_ROLE][%03d ID].[SECONDARY_ROLE].[POOL].[domain]
