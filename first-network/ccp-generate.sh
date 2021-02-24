#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
		-e "s/\${IORG}/$7/" \
        ccp-template.json 
}

function yaml_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n        /g'
}

ORG=collector
P0PORT=7051
P1PORT=8051
CAPORT=7054
PEERPEM=crypto-config/peerOrganizations/collector.spchain.com/tlsca/tlsca.collector.spchain.com-cert.pem
CAPEM=crypto-config/peerOrganizations/collector.spchain.com/ca/ca.collector.spchain.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM ${ORG^})" > connection-collector.json
echo "$(jq '.client.organization |= "Collector" | .organizations.Collector.mspid |= "CollectorMSP"' connection-collector.json)" > connection-collector.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM)" > connection-collector.yaml

ORG=creator
P0PORT=9051
P1PORT=10051
CAPORT=8054
PEERPEM=crypto-config/peerOrganizations/creator.spchain.com/tlsca/tlsca.creator.spchain.com-cert.pem
CAPEM=crypto-config/peerOrganizations/creator.spchain.com/ca/ca.creator.spchain.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM ${ORG^})" > connection-creator.json
echo "$(jq '.client.organization |= "Creator" | .organizations.Creator.mspid |= "CreatorMSP"' connection-creator.json)" > connection-creator.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM)" > connection-creator.yaml

ORG=md
P0PORT=11051
P1PORT=12051
CAPORT=9054
PEERPEM=crypto-config/peerOrganizations/md.spchain.com/tlsca/tlsca.md.spchain.com-cert.pem
CAPEM=crypto-config/peerOrganizations/md.spchain.com/ca/ca.md.spchain.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM ${ORG^})" > connection-md.json
echo "$(jq '.client.organization |= "Md" | .organizations.Md.mspid |= "MdMSP"' connection-md.json)" > connection-md.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM)" > connection-md.yaml

ORG=gallery
P0PORT=13051
P1PORT=14051
CAPORT=10054
PEERPEM=crypto-config/peerOrganizations/gallery.spchain.com/tlsca/tlsca.gallery.spchain.com-cert.pem
CAPEM=crypto-config/peerOrganizations/gallery.spchain.com/ca/ca.gallery.spchain.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM ${ORG^})" > connection-gallery.json
echo "$(jq '.client.organization |= "Gallery" | .organizations.Gallery.mspid |= "GalleryMSP"' connection-gallery.json)" > connection-gallery.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $CAPORT $PEERPEM $CAPEM)" > connection-gallery.yaml
