"use client";

import { motion, type Variants } from "framer-motion";
import { Children, isValidElement, type ReactNode } from "react";

interface StaggerContainerProps {
  readonly children: ReactNode;
  readonly className?: string;
  readonly delay?: number;
}

const containerVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.05,
    },
  },
};

const itemVariants: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.3,
      ease: [0.4, 0, 0.2, 1],
    },
  },
};

export function StaggerContainer({
  children,
  className,
  delay = 0,
}: StaggerContainerProps) {
  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className={className}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {Children.map(children, (child, index) => {
        if (isValidElement(child)) {
          // If child is already a motion component, don't wrap it
          if (child.type === motion.div || child.type === motion.span || child.type === motion.section) {
            return child;
          }
          return (
            <motion.div key={child.key ?? index} variants={itemVariants}>
              {child}
            </motion.div>
          );
        }
        return (
          <motion.div key={index} variants={itemVariants}>
            {child}
          </motion.div>
        );
      })}
    </motion.div>
  );
}
