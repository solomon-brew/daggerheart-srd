# Converting from PDF

    pipx install --python python3.12 marker-pdf
    marker_single \
        --disable_ocr \
        --use_llm \
        --llm_service marker.services.openai.OpenAIService \
        --config_json marker_config.json \
        --openai_api_key "$OPENAI_API_KEY" \
        --openai_model gpt-5.1 \
        --timeout 180 \
        --max_concurrency 1 \
        --output_format markdown \
        --output_dir out \
        DH-SRD-2025-09-09.pdf
    npx --yes markdownlint-cli SRD.md

And then a lot of manual markdown cleanup...
