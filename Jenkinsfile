node {
    script {
//         mysql_addr = '127.0.0.1' // service cluster ip
//         redis_addr = '127.0.0.1' // service cluster ip
        user_addr = '127.0.0.1:30036' // nodeIp : port
    }
//     stage('clone code from github') {
//         echo "first stage: clone code"
//         git url: "https://github.com/liubo51617/user.git"
//         script {
//             commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
//         }
//     }
//     stage('build image') {
//         echo "second stage: build docker image"
//         sh "docker build -t alexliu51617/user:${commit_id} ."
//     }
//     stage('push image') {
//         echo "third stage: push docker image to registry"
//         sh "docker login -u alexliu51617 -p liubo51617"
//         sh "docker push alexliu51617/user:${commit_id}"
//     }
//     stage('deploy to Kubernetes') {
//         echo "forth stage: deploy to Kubernetes"
//         sh "sed -i 's/<COMMIT_ID_TAG>/${commit_id}/' user.yaml"
//         sh "sed -i 's/<MYSQL_ADDR_TAG>/${mysql_addr}/' user.yaml"
//         sh "sed -i 's/<REDIS_ADDR_TAG>/${redis_addr}/' user.yaml"
//         sh "kubectl apply -f user.yaml"
//     }
    stage('http test') {
        echo "fifth stage: http test"
        sh "cd user/service && go test  -args ${user_addr}"
    }
}