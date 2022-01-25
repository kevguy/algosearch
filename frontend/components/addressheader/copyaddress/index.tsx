import { Button } from "@mui/material";
import React, { useState } from "react";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import CheckIcon from "@mui/icons-material/Check";

const CopyAddress = ({
  address,
  className,
}: {
  address: string;
  className: string;
}) => {
  const [copied, setCopied] = useState(false);

  const copyAddress = () => {
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(address);
      setCopied(true);
      setTimeout(() => {
        setCopied(false);
      }, 1000);
    }
  };

  return (
    <Button
      className={className}
      onClick={copyAddress}
      size="small"
      aria-label="copy"
    >
      {copied ? (
        <CheckIcon fontSize="small" />
      ) : (
        <ContentCopyIcon fontSize="small" />
      )}
    </Button>
  );
};

export default CopyAddress;
