TIME_STAMP=`date "+%Y-%m-%d_%H:%M"`

zip ./zip/${TIME_STAMP} *.go go.mod go.sum
gsutil cp ./zip/${TIME_STAMP}.zip gs://mlb-line-bot/


NAME="mlb_line_bot-cf"
ENTRY_POINT="Function"
RUNTIME="go116"
TRIGGER="--trigger-http"

gcloud functions deploy ${NAME} --entry-point ${ENTRY_POINT} --runtime ${RUNTIME} ${TRIGGER} --source=gs://mlb-line-bot/${TIME_STAMP}.zip