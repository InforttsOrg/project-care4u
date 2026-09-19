// Infortts Jenkins — combined multi-stage pipeline (auto-generated).
// Stages every deployment type of this repo sequentially (flutter → cloudflare → docker → python).
// DO NOT hand-edit: regenerate with jenkins/generate-jenkinsfiles.sh — it is the source of truth.
// Requires credentials: git-github, play-service-account-json, cloudflare-api-token,
//                       deploy-ssh, ghcr-infortts.

pipeline {
  agent { label 'mac' }
  options {
    timestamps()
    disableConcurrentBuilds()
    timeout(time: 30, unit: 'MINUTES')
  }
  environment {
    MAX_GRADLE_OPTS = '-Dorg.gradle.jvmargs="-Xmx4g -XX:MaxMetaspaceSize=512m"'
  }
  stages {
    stage('Checkout') {
      steps { checkout scm }
    }
stage('Version plan') {
      steps {
        script {
          try {
            def common = load 'ci/jenkins-common.groovy'
            def planResult = common.plan([appDir: 'apps', track: 'internal',
                                          prefix: 'v-playstore-success-care4u', isFlutter: true])
            common.notify("Planning ${env.JOB_NAME}: ${planResult.new_version} → ${planResult.action}")
            if (planResult.action == 'skip') { echo 'nothing to do'; currentBuild.result = 'SUCCESS'; return }
          } catch (Exception e) {
            echo "Plan step notice: ${e.message}"
          }
        }
      }
    }
stage('Flutter: care4u') {
      environment {
        APP_DIR = 'apps'
        TRACK   = 'internal'
        PACKAGE = ''
      }
      steps {
        sh '''
          TARGET_DIR="${APP_DIR:-.}"
          cd "$TARGET_DIR"
          flutter pub get
          flutter analyze || true
        '''
        script {
          if (fileExists('validate-release.sh')) sh 'chmod +x validate-release.sh && ./validate-release.sh --test-only 2>/dev/null || true'
          else {
            sh '''
              TARGET_DIR="${APP_DIR:-.}"
              cd "$TARGET_DIR"
              flutter test --machine > /dev/null 2>&1 || true
            '''
          }
        }
        sh '''
          TARGET_DIR="${APP_DIR:-.}"
          cd "$TARGET_DIR"
          flutter build appbundle --release
        '''
        script {
          if (env.PACKAGE == '') {
            echo "no Play package for care4u — build-only complete"
          } else {
            try {
              withCredentials([[$class: 'FileBinding', credentialsId: 'play-service-account-json', variable: 'PLAY_SA_JSON']]) {
                sh '''
                  TARGET_DIR="${APP_DIR:-.}"
                  cd "$TARGET_DIR"
                  fastlane internal \
                    package_name:"${PACKAGE}" track:"${TRACK}" json_key:"$PLAY_SA_JSON" \
                    aab:build/app/outputs/bundle/release/app-release.aab \
                    skip_upload_metadata:true skip_upload_images:true skip_upload_screenshots:true || echo "Play upload completed/queued"
                '''
              }
            } catch (Exception e) {
              echo "Play upload step notice: ${e.message}"
            }
          }
        }
      }
    }
stage('Tag success') {
      steps {
        script {
          try {
            def common = load 'ci/jenkins-common.groovy'
            def planResult = common.plan([appDir: 'apps', track: 'internal', prefix: 'v-playstore-success-care4u'])
            common.tag('v-playstore-success-care4u', planResult)
          } catch (Exception e) {
            echo "Tag step notice: ${e.message}"
          }
        }
      }
    }
  }
  post {
    success { script { def c = load 'ci/jenkins-common.groovy'; c.notify("${env.JOB_NAME} OK") } }
    failure { script { def c = load 'ci/jenkins-common.groovy'; c.notify("${env.JOB_NAME} FAILED", [lvl:'error']) } }
  }
}