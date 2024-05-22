import json
import requests
import re

# Send a request to fetch the data
headers = {
    'Referer': 'https://pokedex.org/js/worker.js',
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'
}

skim_urls = "https://pokedex.org/assets/skim-monsters-%s.txt"
description_urls = "https://pokedex.org/assets/descriptions-%s.txt"
monster_supplemental_urls = "https://pokedex.org/assets/monsters-supplemental-%s.txt"
monster_moves_urls = "https://pokedex.org/assets/monster-moves-%s.txt"

final = []


def process_data(resp):
    data = []
    ii = 0
    for line in resp.text.strip().split('\n'):
        if line.startswith('{"seq":'):
            # Extract the seq value
            current_seq = int(json.loads(line)['seq'])
        elif not line.startswith('{"version":'):
            # Parse the JSON data and add seq as id for "docs" entries
            json_data = json.loads(line)
            if 'docs' in json_data:
                for doc in json_data['docs']:
                    descriptions = doc.get('descriptions')
                    if descriptions is not None:
                        resource_uri = descriptions[0]["resource_uri"]
                        number = re.search(r'/(\d+)/$', resource_uri).group(1)
                        doc['descriptions'][0]['resource_uri'] = number
                    types = doc.get('types')
                    if types is not None:
                        for t in types:
                            number = re.search(r'/(\d+)/$', t['resource_uri']).group(1)
                            t['resource_uri'] = number
                    data.append(doc)
    return data


for url in range(1, 4):
    resp = requests.get(skim_urls % str(url), headers=headers)
    skim_data = process_data(resp)
    description_data = process_data(requests.get(description_urls % str(url), headers=headers))
    monster_supplemental_data = process_data(requests.get(monster_supplemental_urls % str(url), headers=headers))
    monster_moves_data = process_data(requests.get(monster_moves_urls % str(url), headers=headers))
    for i in range(len(skim_data)):
        skim_data[i].update(monster_supplemental_data[i])
        skim_data[i].update(monster_moves_data[i])
        final.append(skim_data[i])

with open('final_data.json', 'w') as json_file:
    json.dump(final, json_file, indent=4)




