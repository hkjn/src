import json
data=json.loads(open('prices.json').read())
for d in sorted(data['bpi']):
	price_btc = data['bpi'][d]
	satoshis_per_usd = 10.0**8/price_btc
        satoshis = '{0:.2f} sat'.format(satoshis_per_usd)
        if satoshis_per_usd > 10**8:
            satoshis = '{0:.2f} Gsat'.format(satoshis_per_usd/10.0**9)
        elif satoshis_per_usd > 10**5:
            satoshis = '{0:.2f} Msat'.format(satoshis_per_usd/10.0**6)
        print('{}: {}'.format(d, satoshis))
