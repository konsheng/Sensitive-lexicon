pipeline {
  agent {label 'dockeragent'}
  // 构建逻辑已迁移到 Dockerfile，Jenkins 不再进行本地 go build

  environment {
    GO111MODULE = 'on'        // 开启 Modules 模式
    CGO_ENABLED = '0'
    APP_NAME = 'sensitive-lexicon'
    REGISTRY = 'crpi-vqe38j3xeblrq0n4.cn-hangzhou.personal.cr.aliyuncs.com/go-mctown'
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    // 使用 Dockerfile 完成编译与打包，仅保留镜像构建与推送
    stage('Docker Build & Push') {
      steps {
        withCredentials([usernamePassword(
          credentialsId: 'aliyun-docker-login',
          usernameVariable: 'DOCKER_USERNAME',
          passwordVariable: 'DOCKER_PASSWORD'
        )]) {
          sh """
            echo "\$DOCKER_PASSWORD" | docker login --username \$DOCKER_USERNAME --password-stdin ${env.REGISTRY.split('/')[0]}
          """
        }

        script {
          def imageTag = "${env.REGISTRY}/${env.APP_NAME}:${env.BUILD_NUMBER}"
          def latestTag = "${env.REGISTRY}/${env.APP_NAME}:latest"

          sh """
            ls -l
            docker build -t ${imageTag} --network=host .
            docker tag ${imageTag} ${latestTag}
            docker push ${imageTag}
            docker push ${latestTag}
          """
        }
      }
    }

stage('Deploy All Compose Projects') {
      parallel {
        stage('Deploy compose1') {
  agent {label 'dockeragent'}
          steps {

              checkout scm
                              sh """
                pwd
                ls -l
                """
            dir('deploy/compose') {
              script {
                withCredentials([usernamePassword(
                  credentialsId: 'aliyun-docker-login',
                  usernameVariable: 'DOCKER_USERNAME',
                  passwordVariable: 'DOCKER_PASSWORD'
                )]) {
                  sh """
                    echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin ${env.REGISTRY.split('/')[0]}
                  """
                }
                sh """
                pwd
                ls -l
                docker compose -f docker-compose.yml down || true
                docker compose -f docker-compose.yml pull
                docker compose -f docker-compose.yml up -d --remove-orphans
                """
              }
            }
          }
        }

      }
    }


  }

  post {
    always {
      cleanWs()
    }
    success {
      echo "✅ 构建成功！"
    }
    failure {
      echo "🔥 构建失败，请检查日志。"
    }
  }
}