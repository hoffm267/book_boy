FROM python:3.12-alpine AS builder

WORKDIR /app

RUN apk add --no-cache --no-scripts gcc musl-dev libffi-dev && \
    apk fix --no-scripts

COPY requirements.txt .

RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

FROM python:3.12-alpine AS development

WORKDIR /app

COPY --from=builder /opt/venv /opt/venv

COPY events.py .
COPY consumer.py .

ENV PATH="/opt/venv/bin:$PATH"
ENV PYTHONUNBUFFERED=1

EXPOSE 8000

CMD ["python", "consumer.py"]

FROM python:3.12-alpine AS production

WORKDIR /app

COPY --from=builder /opt/venv /opt/venv

COPY events.py .
COPY consumer.py .

ENV PATH="/opt/venv/bin:$PATH"
ENV PYTHONUNBUFFERED=1

RUN adduser -D -u 1000 appuser 2>/dev/null || true && \
    chown -R appuser:appuser /app 2>/dev/null || true
USER appuser

EXPOSE 8000

CMD ["python", "consumer.py"]
