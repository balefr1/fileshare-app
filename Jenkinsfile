node {
    stage ("SCM checkout"){
        git credentialsId: 'github', url: 'https://github.com/balefr1/go-api'
    }
    stage ("Docker build"){
        sh 'docker build -t --rm eu.gcr.io/homeapi-283920/go-api:latest .'
    }

    withCredentials([string(credentialsId: 'gcr-access-sa', variable: 'gcraccessa')]) {
        sh '''
             docker login -u _json_key -p "${gcraccessa}" https://eu.gcr.io/
         '''
    }

    stage ("Docker push to gcr repo"){
        sh 'docker push eu.gcr.io/homeapi-283920/go-api:latest'
    }

    stage ("Docker clean"){
        sh 'docker image prune -f'
    }
} 