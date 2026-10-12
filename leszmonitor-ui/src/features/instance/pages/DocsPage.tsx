import { useQuery } from "@tanstack/react-query";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { DocsApi } from "@/features/instance/docs-api";
import { readToken } from "@/features/auth/lib/token";
import { QUERY_KEYS, SERVER_API_URL } from "@/lib/consts";

const serverUrl = new URL(SERVER_API_URL);

export function DocsPage() {
  const { data: spec } = useQuery({
    queryFn: DocsApi.getSpec,
    queryKey: [QUERY_KEYS.API_SPEC],
  });

  return (
    <PageContainer>
      <TypographyH1>Documentation</TypographyH1>
      {spec && (
        <SwaggerUI
          spec={{
            ...spec,
            host: serverUrl.host,
            schemes: [serverUrl.protocol.replace(":", "")],
            basePath: serverUrl.pathname,
          }}
          requestInterceptor={async (req) => {
            req.headers.Authorization = `Bearer ${await readToken()}`;
            return req;
          }}
        />
      )}
    </PageContainer>
  );
}
