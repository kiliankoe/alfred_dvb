//
//  AppDelegate.swift
//  ScheduleNotification
//
//  Created by Kilian Költzsch on 24/08/16.
//  Copyright © 2016 Kilian Koeltzsch. All rights reserved.
//

import Cocoa

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate {

    func applicationDidFinishLaunching(_ aNotification: Notification) {

        defer {
            NSApplication.shared().terminate(self)
        }

        let argv = ProcessInfo.processInfo.arguments

        guard argv.count >= 4 else {
            print("You have to supply at least 3 params: deliveryTime, title, informativeText")
            return
        }

        let deliveryTimeString = argv[1]
        let title = argv[2]
        let informativeText = argv[3]

        guard let deliveryTime = Int(deliveryTimeString) else {
            print("Couldn't convert deliveryTime into Integer value.")
            return
        }

        var deliveryDate = Date()
        deliveryDate.addTimeInterval(Double(60) * Double(deliveryTime))

        let notification = NSUserNotification()
        notification.title = title
        notification.informativeText = informativeText

        notification.deliveryDate = deliveryDate

        NSUserNotificationCenter.default.scheduleNotification(notification)
    }
}
