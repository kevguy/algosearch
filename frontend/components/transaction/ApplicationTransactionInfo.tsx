import React from "react";
import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import styles from "../../pages/tx/TransactionDetails.module.scss";
import blockStyles from "../../pages/block/Block.module.scss";

const ApplicationTransactionInfo = ({
  transaction,
}: {
  transaction: TransactionResponse;
}) => (
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
            <td>{transaction["application-transaction"]["application-id"]}</td>
          </tr>
          <tr>
            <td className={styles["valign-top-identifier"]}>Accounts</td>
            <td className={styles["multiline-details"]}>
              {transaction["application-transaction"].accounts.map((ac) => (
                <Link href={`/address/${ac}`} key={ac}>
                  {ac}
                </Link>
              ))}
            </td>
          </tr>
          {transaction["application-transaction"]["application-args"].length >
            0 && (
            <tr>
              <td>Arguments (base64)</td>
              <td className={styles["multiline-details"]}>
                {transaction["application-transaction"]["application-args"].map(
                  (appArg) => (
                    <p key={appArg}>{appArg}</p>
                  )
                )}
              </td>
            </tr>
          )}
          {transaction["application-transaction"]["foreign-apps"].length >
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
          {transaction["application-transaction"]["foreign-assets"].length >
            0 && (
            <tr>
              <td>Foreign Assets</td>
              <td className={styles["multiline-details"]}>
                {transaction["application-transaction"]["foreign-assets"].map(
                  (asa) => (
                    <p key={asa}>{asa}</p>
                  )
                )}
              </td>
            </tr>
          )}
          <tr>
            <td>On Completion</td>
            <td>{transaction["application-transaction"]["on-completion"]}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
);

export default ApplicationTransactionInfo;
