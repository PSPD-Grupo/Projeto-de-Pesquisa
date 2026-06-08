package util;

import java.util.Map;

/**
 * Utilitário mínimo de JSON — evita dependência de Jackson/Gson
 * para manter o projeto enxuto, assim como o original não usa
 * frameworks além do necessário.
 */
public class JsonUtil {

    private JsonUtil() {}

    /** Serializa um Map de chave→valor em JSON simples. */
    public static String toJson(Map<String, Object> fields) {

        StringBuilder sb = new StringBuilder("{");
        boolean first = true;

        for (Map.Entry<String, Object> entry : fields.entrySet()) {

            if (!first) sb.append(",");
            first = false;

            sb.append("\"").append(entry.getKey()).append("\":");

            Object val = entry.getValue();
            if (val instanceof String) {
                sb.append("\"").append(val).append("\"");
            } else {
                sb.append(val);
            }
        }

        sb.append("}");
        return sb.toString();
    }

    /**
     * Extrai o valor de uma chave string em um JSON plano do tipo:
     * {"chave":"valor"} ou {"chave":123}
     */
    public static String extractString(String json, String key) {

        // Tenta com aspas: "key":"value"
        String quotedPattern = "\"" + key + "\":\"";
        int idx = json.indexOf(quotedPattern);

        if (idx >= 0) {
            int start = idx + quotedPattern.length();
            int end   = json.indexOf("\"", start);
            return end >= 0 ? json.substring(start, end) : null;
        }

        // Tenta sem aspas (número): "key":123
        String plainPattern = "\"" + key + "\":";
        idx = json.indexOf(plainPattern);

        if (idx >= 0) {
            int start = idx + plainPattern.length();
            int end   = json.indexOf("}", start);
            String raw = end >= 0 ? json.substring(start, end) : null;
            return raw != null ? raw.replaceAll("[^0-9\\-]", "") : null;
        }

        return null;
    }
}