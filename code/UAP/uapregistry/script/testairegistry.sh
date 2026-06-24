#!/bin/bash

unset http_proxy https_proxy no_proxy all_proxy

function print() {
    echo ""
    echo ""
    echo "-----------------------------------------------"
    echo "----   ${1}"
    echo "-----------------------------------------------"
    echo ""
}

function createMultiNodes() {
    print "创建多个node"
    nodes='[
    {"Labels":["player", "student"],"Props":{"name": "2", "age": 15}},
    {"Labels":["student"],"Props":{"name": "3", "age": 16}},
    {"Labels":["student"],"Props":{"name": "4", "age": 16}},
    {"Labels":["student"],"Props":{"name": "5", "age": 16}},
    {"Labels":["student"],"Props":{"name": "6", "age": 16}},
    {"Labels":["student"],"Props":{"name": "6", "age": 16}},
    {"Labels":["student"],"Props":{"name": "8", "age": 16}},
    {"Labels":["student"],"Props":{"name": "9", "age": 16}},
    {"Labels":["student"],"Props":{"name": "10", "age": 16}}
    ]
    '
    curl -X POST -w "\n%{http_code}" 'http://localhost:8080/knowledgegraph/nodes/bulk' -d "${nodes}"

}

function createNode() {
    print "创建单个node"
    node=$(curl -s -X POST 'http://localhost:8080/knowledgegraph/nodes' -d '{"Labels":["player", "student"],"Props":{"name": "1", "age": 15}}')

    echo "${node}"

    nodeid=$(echo "$node" | jq -r '.ElementId')
    echo "nodeid is $nodeid"
}


function updateNode() {
    node='{"Labels":["player", "student"],"Props":{"name": "Aa", "age": 15}}'

    curl -w "\n%{http_code}" -X PUT http://localhost:8080/knowledgegraph/nodes/"${1}" -d "${node}"
}


function createRelationship(){
    print "创建单个relationship"
    node1=$(curl -s -X GET "http://127.0.0.1:8080/knowledgegraph/nodes?page=1&limit=1")
    echo "$node1"
    node1id=$(echo "$node1"| jq -r '.[0].ElementId')
    node2=$(curl -s -X GET "http://127.0.0.1:8080/knowledgegraph/nodes?page=1&limit=1")
    echo "$node2"
    node2id=$(echo "$node2"| jq -r '.[0].ElementId')

    relationship='
    {"StartElementId":"'"${node1id}"'","EndElementId":"'"${node2id}"'","Type":"classmate","Props":{"grade":"1", "class":9}}
    '
    curl -s -w "\n%{http_code}" -X POST 'http://localhost:8080/knowledgegraph/relationships' -d "${relationship}"
}

function createRelationships() {
    print "创建多个relationship"

    node3=$(curl -s -X GET "http://127.0.0.1:8080/knowledgegraph/nodes?page=2&limit=1")
    echo "$node3"
    node3id=$(echo "$node3"| jq -r '.[0].ElementId')

    node4=$(curl -s -X GET "http://127.0.0.1:8080/knowledgegraph/nodes?page=3&limit=1")
    echo "$node4"
    node4id=$(echo "$node4"| jq -r '.[0].ElementId')


    relationships='[
    {"StartElementId":"'"${node1id}"'","EndElementId":"'"${node2id}"'","Type":"classmate","Props":{"grade":"1", "class":9}}, 
    {"StartElementId":"'"${node1id}"'","EndElementId":"'"${node3id}"'","Type":"classmate","Props":{"grade":"1", "class":9}}, 
    {"StartElementId":"'"${node2id}"'","EndElementId":"'"${node4id}"'","Type":"roommate","Props":{"roomnum":"201", "building":"aaaaa"}}
    ]
    '
    curl -s -w "\n%{http_code}" -X POST 'http://localhost:8080/knowledgegraph/relationships/bulk' -d "${relationships}"
}


function updateRelationship() {
    print "更新relationship"
    relationship1=$(curl -s "http://localhost:8080/knowledgegraph/relationships?page=1&limit=1")
    echo "old relationship: "
    echo "${relationship1}"

    relationship1id=$(echo "${relationship1}" |  jq -r '.[0].ElementId')

    echo "update"
    curl -w "\n%{http_code}" -X PUT "http://localhost:8080/knowledgegraph/relationships/${relationship1id}" -d '{"updated":"true"}'
}


function deleteAllNodes() {
    nodeExist=true

    while $nodeExist; do
        node=$(curl -s -X GET "http://127.0.0.1:8080/knowledgegraph/nodes?page=1&limit=1")
        if [ "$node" = "null" ]; then
            nodeExist=false
        else
            echo "$node"
            nodeid=$(echo "$node" | jq -r '.[0].ElementId')
            response=$(curl -s -w "\n%{http_code}" -X DELETE "http://127.0.0.1:8080/knowledgegraph/nodes/${nodeid}")
            echo "${response}"
            if [[ "${response}" == *relationships* ]]; then
                echo "force delete"
                curl -s -w "\n%{http_code}" -X DELETE "http://127.0.0.1:8080/knowledgegraph/nodes/${nodeid}?force=true"
            fi
        fi

    done
    echo "已经删除所有节点"
}

nodeid="empty"
node1id=""
node2id=""
relationship1id=""

function main() {

    print "导出graph"
    curl -w "\n%{http_code}" -X GET http://127.0.0.1:8080/knowledgegraph/graph

    createNode

    createMultiNodes

    print "更新node"
    updateNode "${nodeid}"

    print "查询单个node"
    curl -s -w "\n%{http_code}" -X GET  http://localhost:8080/knowledgegraph/nodes/"${nodeid}"

    print "查询不存在的单个node"
    curl -s -w "\n%{http_code}" -X GET  http://localhost:8080/knowledgegraph/nodes/lvbotest

    print "获取节点列表分页查询"
    curl -X GET -w "\n%{http_code}" "http://127.0.0.1:8080/knowledgegraph/nodes?page=1&limit=2"

    print "获取节点列表"
    curl -X GET -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/nodes

    printf "\n*************************************************************************************"
    echo "***** relationship 测试"
    echo "*************************************************************************************"

    createRelationship

    createRelationships

    updateRelationship

    print "查询单个relationship"
    curl -s -w "\n%{http_code}" -X GET  "http://localhost:8080/knowledgegraph/relationships/${relationship1id}"

    print "查询单个不存在的relationship"
    curl -s -w "\n%{http_code}" -X GET  "http://localhost:8080/knowledgegraph/relationships/lvbotest"


    print "获取关系列表分页查询"
    curl -X GET -w "\n%{http_code}" "http://127.0.0.1:8080/knowledgegraph/relationships?page=1&limit=2"

    print "获取关系列表"
    curl -X GET -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/relationships


    print "导出graph"
    graph=$(curl -s -X GET http://127.0.0.1:8080/knowledgegraph/graph)
    echo "${graph}"

    graphData=$(echo "$graph" | jq -r .data)

    print "开始删除节点"
    deleteAllNodes


    print "导出graph"
    curl -X GET -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/graph


    print "导入graph"
    curl -X POST -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/graph -d "${graphData}"


    print "导出graph"
    curl -X GET -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/graph


    print "开始删除节点"
    deleteAllNodes

    print "导出graph"
    curl -X GET -w "\n%{http_code}" http://127.0.0.1:8080/knowledgegraph/graph

    echo ""

    }

main