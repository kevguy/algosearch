import React, { useEffect, useState } from "react";
import Link from "next/link";
import algosdk from "algosdk";
import hljs from "highlight.js";

import { TransactionResponse } from "../../types/apiResponseTypes";
import styles from "../../pages/tx/TransactionDetails.module.scss";
import blockStyles from "../../pages/block/Block.module.scss";
import { getAppTEAL, getClearStateTEAL } from "../../utils/api";
import { prettyPrintTEAL } from "../../utils/stringUtils";

const ApplicationTransactionInfo = ({
  transaction,
}: {
  transaction: TransactionResponse;
}) => {
  const [disassembledApp, setDisassembledApp] = useState<string>();
  const [disassembledClearStateProgram, setDisassembledClearStateProgram] =
    useState<string>();

  useEffect(() => {
    if (transaction && transaction["application-transaction"]) {
      if (transaction["application-transaction"]["approval-program"]) {
        getAppTEAL(transaction)
          .then((result) => {
            if (
              result &&
              result.txns &&
              result.txns[0] &&
              result.txns[0].disassembly
            ) {
              const disassembledResult = prettyPrintTEAL(
                result.txns[0].disassembly
              );
              setDisassembledApp(disassembledResult);
            }
          })
          .catch((error) => {
            console.error("App disassembly error: ", error);
          });
      }
      if (transaction["application-transaction"]["clear-state-program"]) {
        getClearStateTEAL(transaction)
          .then((result) => {
            if (
              result &&
              result.txns &&
              result.txns[0] &&
              result.txns[0].disassembly
            ) {
              const disassembledResult = prettyPrintTEAL(
                result.txns[0].disassembly
              );
              setDisassembledClearStateProgram(disassembledResult);
            }
          })
          .catch((error) => {
            console.error("App disassembly error: ", error);
          });
      }
    }
  }, [transaction]);

  return (
    <div>
      <h4>Application Transaction Information</h4>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <thead>
            <tr>
              <th>Identifier</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Application ID</td>
              <td>
                {transaction["application-transaction"]["application-id"]}
              </td>
            </tr>
            <tr>
              <td className={styles["valign-top-identifier"]}>Accounts</td>
              <td className={styles["multiline-details"]}>
                {transaction["application-transaction"].accounts &&
                transaction["application-transaction"].accounts.length > 0
                  ? transaction["application-transaction"].accounts.map(
                      (ac) => (
                        <Link href={`/address/${ac}`} key={ac}>
                          {ac}
                        </Link>
                      )
                    )
                  : "N/A"}
              </td>
            </tr>
            {transaction["application-transaction"]["application-args"] &&
              transaction["application-transaction"]["application-args"]
                .length > 0 && (
                <tr className={styles["valign-top-identifier"]}>
                  <td>Arguments (base64)</td>
                  <td className={styles["multiline-details"]}>
                    <div className={styles["inner-table-wrapper"]}>
                      <table className={styles["inner-table"]}>
                        <thead>
                          <tr>
                            <td>base64</td>
                            <td>ascii</td>
                            <td>uint</td>
                          </tr>
                        </thead>
                        <tbody>
                          {transaction["application-transaction"][
                            "application-args"
                          ].map((appArg) => (
                            <tr key={appArg}>
                              <td>{appArg}</td>
                              <td>
                                {Buffer.from(appArg, "base64").toString(
                                  "ascii"
                                )}
                              </td>
                              <td>
                                {Buffer.from(appArg, "base64").length <= 8
                                  ? algosdk.decodeUint64(
                                      Buffer.from(appArg, "base64"),
                                      "mixed"
                                    )
                                  : "N/A"}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </td>
                </tr>
              )}
            {transaction["application-transaction"]["foreign-apps"] &&
              transaction["application-transaction"]["foreign-apps"].length >
                0 && (
                <tr>
                  <td>Foreign Apps</td>
                  <td className={styles["multiline-details"]}>
                    {transaction["application-transaction"]["foreign-apps"].map(
                      (app) => (
                        <p key={app}>{app}</p>
                      )
                    )}
                  </td>
                </tr>
              )}
            {transaction["application-transaction"]["foreign-assets"] &&
              transaction["application-transaction"]["foreign-assets"].length >
                0 && (
                <tr>
                  <td>Foreign Assets</td>
                  <td className={styles["multiline-details"]}>
                    {transaction["application-transaction"][
                      "foreign-assets"
                    ].map((asa) => (
                      <p key={asa}>{asa}</p>
                    ))}
                  </td>
                </tr>
              )}
            <tr>
              <td>On Completion</td>
              <td>{transaction["application-transaction"]["on-completion"]}</td>
            </tr>
            {transaction["application-transaction"]["approval-program"] &&
              disassembledApp && (
                <tr>
                  <td className={styles["valign-top-identifier"]}>
                    Approval Program
                  </td>
                  <td>
                    <pre className={`${styles["teal-box"]} hljs`}>
                      <code
                        className="language-lua"
                        dangerouslySetInnerHTML={{
                          __html: hljs.highlight(disassembledApp, {
                            language: "lua",
                          }).value,
                        }}
                      ></code>
                    </pre>
                  </td>
                </tr>
              )}
            {transaction["application-transaction"]["clear-state-program"] &&
              disassembledClearStateProgram && (
                <tr>
                  <td className={styles["valign-top-identifier"]}>
                    Clear State Program
                  </td>
                  <td>
                    <pre className={`${styles["teal-box"]} hljs`}>
                      <code
                        className="language-lua"
                        dangerouslySetInnerHTML={{
                          __html: hljs.highlight(
                            disassembledClearStateProgram,
                            {
                              language: "lua",
                            }
                          ).value,
                        }}
                      ></code>
                    </pre>
                  </td>
                </tr>
              )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
export default ApplicationTransactionInfo;
