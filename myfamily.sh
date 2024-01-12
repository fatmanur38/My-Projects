
curl -s https://talent.uniworkhub.com/assets/superhero/all.json | jq --arg HERO_ID "$HERO_ID" '.[] | select(.id == ($HERO_ID | tonumber)).connections.relatives' | tr -d '"' 
