curl  https://talent.uniworkhub.com/assets/superhero/all.json | jq '.[] | select(.id == 170)| 
    {
    name,"powerstats": .powerstats.power,"gender":.appearance.gender
    }'